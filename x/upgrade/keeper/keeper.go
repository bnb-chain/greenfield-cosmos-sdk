package keeper

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/cosmos/cosmos-sdk/types/module"
	xp "github.com/cosmos/cosmos-sdk/x/upgrade/exported"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type Keeper struct {
	homePath           string                              // root directory of app config
	storeKey           storetypes.StoreKey                 // key to access x/upgrade store
	cdc                codec.BinaryCodec                   // App-wide binary codec
	upgradeHandlers    map[string]types.UpgradeHandler     // map of plan name to upgrade handler
	upgradeInitializer map[string]types.UpgradeInitializer // map of plan name to upgrade initializer
	upgradeConfig      *types.UpgradeConfig                // upgrade config for upcoming upgrade plan or upgraded plan
	versionSetter      xp.ProtocolVersionSetter            // implements setting the protocol version field on BaseApp
	initVersionMap     module.VersionMap
}

// NewKeeper constructs an upgrade Keeper which requires the following arguments:
// skipUpgradeHeights - map of heights to skip an upgrade
// storeKey - a store key with which to access upgrade's store
// cdc - the app-wide binary codec
// homePath - root directory of the application's config
// vs - the interface implemented by baseapp which allows setting baseapp's protocol version field
func NewKeeper(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, homePath string, vs xp.ProtocolVersionSetter) *Keeper {
	return &Keeper{
		homePath:           homePath,
		storeKey:           storeKey,
		cdc:                cdc,
		upgradeHandlers:    make(map[string]types.UpgradeHandler),
		upgradeInitializer: make(map[string]types.UpgradeInitializer),
		upgradeConfig:      types.NewUpgradeConfig(),
		versionSetter:      vs,
	}
}

// SetVersionSetter sets the interface implemented by baseapp which allows setting baseapp's protocol version field
func (k *Keeper) SetVersionSetter(vs xp.ProtocolVersionSetter) {
	k.versionSetter = vs
}

// GetVersionSetter gets the protocol version field of baseapp
func (k *Keeper) GetVersionSetter() xp.ProtocolVersionSetter {
	return k.versionSetter
}

// SetInitVersionMap sets the initial version map.
// This is only used in app wiring and should not be used in any other context.
func (k *Keeper) SetInitVersionMap(vm module.VersionMap) {
	k.initVersionMap = vm
}

// GetInitVersionMap gets the initial version map
// This is only used in upgrade InitGenesis and should not be used in any other context.
func (k *Keeper) GetInitVersionMap() module.VersionMap {
	return k.initVersionMap
}

// SetUpgradeHandler sets an UpgradeHandler for the upgrade specified by name. This handler will be called when the upgrade
// with this name is applied. In order for an upgrade with the given name to proceed, a handler for this upgrade
// must be set even if it is a no-op function.
func (k Keeper) SetUpgradeHandler(name string, upgradeHandler types.UpgradeHandler) {
	k.upgradeHandlers[name] = upgradeHandler
}

// SetUpgradeInitializer sets an UpgradeInitializer for the upgrade specified by name. This initializer will be called when the program restart after upgrade
// with this name is applied. In order for an upgrade with the given name to proceed, a initializer for this upgrade
// must be set even if it is a no-op function.
func (k Keeper) SetUpgradeInitializer(name string, upgradeInitializer types.UpgradeInitializer) {
	k.upgradeInitializer[name] = upgradeInitializer
}

// SetModuleVersionMap saves a given version map to state
func (k Keeper) SetModuleVersionMap(ctx sdk.Context, vm module.VersionMap) {
	if len(vm) > 0 {
		store := ctx.KVStore(k.storeKey)
		versionStore := prefix.NewStore(store, []byte{types.VersionMapByte})
		// Even though the underlying store (cachekv) store is sorted, we still
		// prefer a deterministic iteration order of the map, to avoid undesired
		// surprises if we ever change stores.
		sortedModNames := make([]string, 0, len(vm))

		for key := range vm {
			sortedModNames = append(sortedModNames, key)
		}
		sort.Strings(sortedModNames)

		for _, modName := range sortedModNames {
			ver := vm[modName]
			nameBytes := []byte(modName)
			verBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(verBytes, ver)
			versionStore.Set(nameBytes, verBytes)
		}
	}
}

// GetModuleVersionMap returns a map of key module name and value module consensus version
// as defined in ADR-041.
func (k Keeper) GetModuleVersionMap(ctx sdk.Context) module.VersionMap {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, []byte{types.VersionMapByte})

	vm := make(module.VersionMap)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		moduleBytes := it.Key()
		// first byte is prefix key, so we remove it here
		name := string(moduleBytes[1:])
		moduleVersion := binary.BigEndian.Uint64(it.Value())
		vm[name] = moduleVersion
	}

	return vm
}

// GetModuleVersions gets a slice of module consensus versions
func (k Keeper) GetModuleVersions(ctx sdk.Context) []*types.ModuleVersion {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, []byte{types.VersionMapByte})
	defer it.Close()

	mv := make([]*types.ModuleVersion, 0)
	for ; it.Valid(); it.Next() {
		moduleBytes := it.Key()
		name := string(moduleBytes[1:])
		moduleVersion := binary.BigEndian.Uint64(it.Value())
		mv = append(mv, &types.ModuleVersion{
			Name:    name,
			Version: moduleVersion,
		})
	}
	return mv
}

// getModuleVersion gets the version for a given module, and returns true if it exists, false otherwise
func (k Keeper) getModuleVersion(ctx sdk.Context, name string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, []byte{types.VersionMapByte})
	defer it.Close()

	for ; it.Valid(); it.Next() {
		moduleName := string(it.Key()[1:])
		if moduleName == name {
			version := binary.BigEndian.Uint64(it.Value())
			return version, true
		}
	}
	return 0, false
}

// ScheduleUpgrade schedules an upgrade based on the specified plan.
// If there is another Plan already scheduled, it will cancel and overwrite it.
// ScheduleUpgrade will also write the upgraded IBC ClientState to the upgraded client
// path if it is specified in the plan.
func (k Keeper) ScheduleUpgrade(ctx sdk.Context, plan types.Plan) error {
	if err := plan.ValidateBasic(); err != nil {
		return err
	}

	// NOTE: allow for the possibility of chains to schedule upgrades in begin block of the same block
	// as a strategy for emergency hard fork recoveries
	if plan.Height < ctx.BlockHeight() {
		return types.ErrUpgradeScheduled
	}

	if k.GetDoneHeight(ctx, plan.Name) != 0 {
		return types.ErrUpgradeCompleted
	}

	// clear any old IBC state stored by previous plan
	oldPlan, found := k.GetUpgradePlan(ctx)
	if found {
		for _, plan := range oldPlan {
			k.ClearIBCState(ctx, plan.Height)
		}
	}

	k.upgradeConfig.SetPlan(&plan)

	return nil
}

// SetUpgradedConsensusState sets the expected upgraded consensus state for the next version of this chain
// using the last height committed on this chain.
func (k Keeper) SetUpgradedConsensusState(ctx sdk.Context, planHeight int64, bz []byte) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.UpgradedConsStateKey(planHeight), bz)
	return nil
}

// GetUpgradedConsensusState gets the expected upgraded consensus state for the next version of this chain
func (k Keeper) GetUpgradedConsensusState(ctx sdk.Context, lastHeight int64) ([]byte, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.UpgradedConsStateKey(lastHeight))
	if len(bz) == 0 {
		return nil, false
	}

	return bz, true
}

// GetLastCompletedUpgrade returns the last applied upgrade name and height.
func (k Keeper) GetLastCompletedUpgrade(ctx sdk.Context) (string, int64) {
	iter := sdk.KVStoreReversePrefixIterator(ctx.KVStore(k.storeKey), []byte{types.DoneByte})
	defer iter.Close()

	if iter.Valid() {
		return parseDoneKey(iter.Key())
	}

	return "", 0
}

// parseDoneKey - split upgrade name and height from the done key
func parseDoneKey(key []byte) (string, int64) {
	// 1 byte for the DoneByte + 8 bytes height + at least 1 byte for the name
	kv.AssertKeyAtLeastLength(key, 10)
	height := binary.BigEndian.Uint64(key[1:9])
	return string(key[9:]), int64(height)
}

// encodeDoneKey - concatenate DoneByte, height and upgrade name to form the done key
func encodeDoneKey(name string, height int64) []byte {
	key := make([]byte, 9+len(name)) // 9 = donebyte + uint64 len
	key[0] = types.DoneByte
	binary.BigEndian.PutUint64(key[1:9], uint64(height))
	copy(key[9:], name)
	return key
}

// GetDoneHeight returns the height at which the given upgrade was executed
func (k Keeper) GetDoneHeight(ctx sdk.Context, name string) int64 {
	iter := sdk.KVStorePrefixIterator(ctx.KVStoreWithZeroRead(k.storeKey), []byte{types.DoneByte})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		upgradeName, height := parseDoneKey(iter.Key())
		if upgradeName == name {
			return height
		}
	}
	return 0
}

// ClearIBCState clears any planned IBC state
func (k Keeper) ClearIBCState(ctx sdk.Context, lastHeight int64) {
	// delete IBC client and consensus state from store if this is IBC plan
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.UpgradedConsStateKey(lastHeight))
}

// ClearUpgradePlan clears any schedule upgrade and associated IBC states.
func (k Keeper) ClearUpgradePlan(ctx sdk.Context) {
	// clear IBC states everytime upgrade plan is removed
	planHeight := ctx.BlockHeight()
	oldPlans, found := k.GetUpgradePlan(ctx)
	if found {
		for _, plan := range oldPlans {
			planHeight = plan.Height
			k.ClearIBCState(ctx, plan.Height)
		}
	}

	k.upgradeConfig.Clear(planHeight)
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetUpgradePlan returns the currently scheduled Plan if any, setting havePlan to true if there is a scheduled
// upgrade or false if there is none
func (k Keeper) GetUpgradePlan(ctx sdk.Context) ([]*types.Plan, bool) {
	plans := k.upgradeConfig.GetPlan(ctx.BlockHeight())
	if len(plans) == 0 {
		return nil, false
	}

	nonUpgraded := make([]*types.Plan, 0, len(plans))
	for i := 0; i < len(plans); i++ {
		if !k.IsUpgraded(ctx, plans[i].Name) {
			nonUpgraded = append(nonUpgraded, plans[i])
		}
	}

	if len(nonUpgraded) == 0 {
		return nil, false
	}

	return nonUpgraded, true
}

// setDone marks this upgrade name as being done so the name can't be reused accidentally
func (k Keeper) setDone(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(encodeDoneKey(name, ctx.BlockHeight()), []byte{1})
}

// HasHandler returns true iff there is a handler registered for this name
func (k Keeper) HasHandler(name string) bool {
	_, ok := k.upgradeHandlers[name]
	return ok
}

// ApplyUpgrade will execute the handler associated with the Plan and mark the plan as done.
func (k Keeper) ApplyUpgrade(ctx sdk.Context, plan types.Plan) {
	initializer := k.upgradeInitializer[plan.Name]

	if initializer != nil {
		err := initializer()
		if err != nil {
			ctx.Logger().Error("failed to init upgrade ["+plan.Name+"]", "err", err)
			return
		}
	}

	handler := k.upgradeHandlers[plan.Name]
	if handler == nil {
		ctx.Logger().Error("missing handler to upgrade [" + plan.Name + "]")
		return
	}

	updatedVM, err := handler(ctx, plan, k.GetModuleVersionMap(ctx))
	if err != nil {
		ctx.Logger().Error("failed to upgrade ["+plan.Name+"]", "err", err)
		return
	}
	k.SetModuleVersionMap(ctx, updatedVM)

	// Must clear IBC state after upgrade is applied as it is stored separately from the upgrade plan.
	// This will prevent resubmission of upgrade msg after upgrade is already completed.
	k.setDone(ctx, plan.Name)
}

// DumpUpgradeInfoToDisk writes upgrade information to UpgradeInfoFileName.
func (k Keeper) DumpUpgradeInfoToDisk(height int64, p types.Plan) error {
	upgradeInfoFilePath, err := k.GetUpgradeInfoPath()
	if err != nil {
		return err
	}

	upgradeInfo := types.Plan{
		Name:   p.Name,
		Height: height,
		Info:   p.Info,
	}
	info, err := json.Marshal(upgradeInfo)
	if err != nil {
		return err
	}

	return os.WriteFile(upgradeInfoFilePath, info, 0o600)
}

// GetUpgradeInfoPath returns the upgrade info file path
func (k Keeper) GetUpgradeInfoPath() (string, error) {
	upgradeInfoFileDir := path.Join(k.getHomeDir(), "data")
	if err := os.MkdirAll(upgradeInfoFileDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create directory %q: %w", upgradeInfoFileDir, err)
	}

	return filepath.Join(upgradeInfoFileDir, types.UpgradeInfoFilename), nil
}

// getHomeDir returns the height at which the given upgrade was executed
func (k Keeper) getHomeDir() string {
	return k.homePath
}

// ReadUpgradeInfoFromDisk returns the name and height of the upgrade which is
// written to disk by the old binary when panicking. An error is returned if
// the upgrade path directory cannot be created or if the file exists and
// cannot be read or if the upgrade info fails to unmarshal.
func (k Keeper) ReadUpgradeInfoFromDisk() (types.Plan, error) {
	var upgradeInfo types.Plan

	upgradeInfoPath, err := k.GetUpgradeInfoPath()
	if err != nil {
		return upgradeInfo, err
	}

	data, err := os.ReadFile(upgradeInfoPath)
	if err != nil {
		// if file does not exist, assume there are no upgrades
		if os.IsNotExist(err) {
			return upgradeInfo, nil
		}

		return upgradeInfo, err
	}

	if err := json.Unmarshal(data, &upgradeInfo); err != nil {
		return upgradeInfo, err
	}

	return upgradeInfo, nil
}

// IsUpgraded returns the bool which the given upgrade was executed
func (k Keeper) IsUpgraded(ctx sdk.Context, name string) bool {
	height := k.GetDoneHeight(ctx, name)
	if height == 0 {
		return false
	}

	return height <= ctx.BlockHeight()
}

// InitUpgraded execute the upgrade initializer that the upgrade is already applied.
func (k Keeper) InitUpgraded(ctx sdk.Context) error {
	iter := storetypes.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte{types.DoneByte})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		upgradeName, height := parseDoneKey(iter.Key())
		if height < ctx.BlockHeight() {
			f := k.upgradeInitializer[upgradeName]
			if f == nil {
				continue
			}

			err := f()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// RegisterUpgradePlan registers the upgrade plans form serverconfig.UpgradeConfig
func (keeper *Keeper) RegisterUpgradePlan(chianID string, plans []serverconfig.UpgradeConfig) error {
	c, err := convertUpgradeConfig(chianID, plans)
	if err != nil {
		return err
	}
	keeper.upgradeConfig = c
	return nil
}

// getExistChainConfig returns the exist chain config
func getExistChainConfig(chainID string) *types.UpgradeConfig {
	switch chainID {
	case types.MainnetChainID:
		return types.MainnetConfig
	case types.TestnetChainID:
		return types.TestnetConfig
	default:
		return types.NewUpgradeConfig()
	}
}

// convertUpgradeConfig converts serverconfig.UpgradeConfig to types.UpgradeConfig
func convertUpgradeConfig(chainID string, plans []serverconfig.UpgradeConfig) (*types.UpgradeConfig, error) {
	upgradeConfig := getExistChainConfig(chainID)

	// override by app config
	for _, plan := range plans {
		nPlan := &types.Plan{
			Name:   plan.Name,
			Height: plan.Height,
			Info:   plan.Info,
		}
		if err := nPlan.ValidateBasic(); err != nil {
			return nil, err
		}
		upgradeConfig.SetPlan(nPlan)
	}

	return upgradeConfig, nil
}
