package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"runtime/debug"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crosschaintypes "github.com/cosmos/cosmos-sdk/x/crosschain/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the oracle MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{
		k,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Claim(goCtx context.Context, req *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	logger := k.Logger(ctx)

	// check dest chain id
	if sdk.ChainID(req.DestChainId) != k.CrossChainKeeper.GetSrcChainID() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDestChainId, "dest chain id(%d) should be %d", req.DestChainId, k.CrossChainKeeper.GetSrcChainID())
	}

	// check src chain id
	if !k.CrossChainKeeper.IsDestChainSupported(sdk.ChainID(req.SrcChainId)) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidSrcChainId, "src chain id(%d) is not supported", req.SrcChainId)
	}

	sequence := k.CrossChainKeeper.GetReceiveSequence(ctx, types.RelayPackagesChannelId)
	if sequence != req.Sequence {
		return nil, sdkerrors.Wrapf(types.ErrInvalidReceiveSequence, "current sequence of channel %d is %d", types.RelayPackagesChannelId, sequence)
	}

	relayer, signedRelayers, err := k.CheckClaim(ctx, req)
	if err != nil {
		return nil, err
	}

	packages := types.Packages{}
	err = rlp.DecodeBytes(req.Payload, &packages)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPayload, "decode payload error")
	}

	events := make([]proto.Message, 0, len(packages))
	totalRelayerFee := sdkmath.ZeroInt()
	for idx := range packages {
		pack := packages[idx]

		relayerFee, event, err := k.handlePackage(ctx, &pack, req.SrcChainId, req.DestChainId, req.Timestamp)
		if err != nil {
			logger.Error("process package failed", "channel", pack.ChannelId, "sequence", pack.Sequence, "error", err.Error())
			return nil, err
		}
		logger.Info("process package success", "channel", pack.ChannelId, "sequence", pack.Sequence)

		events = append(events, event)

		totalRelayerFee = totalRelayerFee.Add(relayerFee)

		// increase channel sequence
		k.CrossChainKeeper.IncrReceiveSequence(ctx, pack.ChannelId)
	}

	err = k.distributeReward(ctx, relayer, signedRelayers, totalRelayerFee)
	if err != nil {
		return nil, err
	}

	k.CrossChainKeeper.IncrReceiveSequence(ctx, types.RelayPackagesChannelId)

	err = ctx.EventManager().EmitTypedEvents(events...)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimResponse{}, nil
}

// distributeReward will distribute reward to relayers
func (k Keeper) distributeReward(ctx sdk.Context, relayer sdk.AccAddress, signedRelayers []sdk.AccAddress, relayerFee sdkmath.Int) error {
	if !relayerFee.IsPositive() {
		k.Logger(ctx).Info("total relayer fee is zero")
		return nil
	}

	otherRelayers := make([]sdk.AccAddress, 0, len(signedRelayers))
	for _, signedRelayer := range signedRelayers {
		if !signedRelayer.Equals(relayer) {
			otherRelayers = append(otherRelayers, relayer)
		}
	}

	totalDistributed, otherRelayerReward := sdkmath.ZeroInt(), sdkmath.ZeroInt()

	relayerRewardShare := k.GetRelayerRewardShare(ctx)

	// calculate the reward to distribute to each other relayer
	if len(otherRelayers) > 0 {
		otherRelayerReward = relayerFee.Mul(sdkmath.NewInt(100 - int64(relayerRewardShare))).Mul(relayerFee).Quo(sdkmath.NewInt(100)).Quo(sdkmath.NewInt(int64(len(otherRelayers))))
	}

	bondDenom := k.StakingKeeper.BondDenom(ctx)
	if otherRelayerReward.IsPositive() {
		for _, signedRelayer := range otherRelayers {
			err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx,
				crosschaintypes.ModuleName,
				signedRelayer,
				sdk.Coins{sdk.Coin{Denom: bondDenom, Amount: otherRelayerReward}},
			)
			if err != nil {
				return err
			}
			totalDistributed = totalDistributed.Add(otherRelayerReward)
		}
	}

	remainingReward := relayerFee.Sub(totalDistributed)
	if remainingReward.IsPositive() {
		err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx,
			crosschaintypes.ModuleName,
			relayer,
			sdk.Coins{sdk.Coin{Denom: bondDenom, Amount: remainingReward}},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) handlePackage(
	ctx sdk.Context,
	pack *types.Package,
	srcChainId uint32,
	destChainId uint32,
	timestamp uint64,
) (sdkmath.Int, *types.EventPackageClaim, error) {
	logger := k.Logger(ctx)

	crossChainApp := k.CrossChainKeeper.GetCrossChainApp(pack.ChannelId)
	if crossChainApp == nil {
		return sdkmath.ZeroInt(), nil, sdkerrors.Wrapf(types.ErrChannelNotRegistered, "channel %d not registered", pack.ChannelId)
	}

	sequence := k.CrossChainKeeper.GetReceiveSequence(ctx, pack.ChannelId)
	if sequence != pack.Sequence {
		return sdkmath.ZeroInt(), nil, sdkerrors.Wrapf(types.ErrInvalidReceiveSequence,
			"current sequence of channel %d is %d", pack.ChannelId, sequence)
	}

	packageHeader, err := sdk.DecodePackageHeader(pack.Payload)
	if err != nil {
		return sdkmath.ZeroInt(), nil, sdkerrors.Wrapf(types.ErrInvalidPayloadHeader, "payload header is invalid")
	}

	if packageHeader.Timestamp != timestamp {
		return sdkmath.ZeroInt(), nil, sdkerrors.Wrapf(types.ErrInvalidPayloadHeader,
			"timestamp(%d) is not the same in payload header(%d)", timestamp, packageHeader.Timestamp)
	}

	if !sdk.IsValidCrossChainPackageType(packageHeader.PackageType) {
		return sdkmath.ZeroInt(), nil, sdkerrors.Wrapf(types.ErrInvalidPackageType,
			"package type %d is invalid", packageHeader.PackageType)
	}

	cacheCtx, write := ctx.CacheContext()
	crash, result := executeClaim(cacheCtx, crossChainApp, sequence, pack.Payload, &packageHeader)
	if result.IsOk() {
		write()
	}

	// write ack package
	var sendSequence int64 = -1
	if packageHeader.PackageType == sdk.SynCrossChainPackageType {
		if crash {
			if len(pack.Payload) < sdk.SynPackageHeaderLength {
				logger.Error("found payload without header",
					"channelID", pack.ChannelId, "sequence", pack.Sequence, "payload", hex.EncodeToString(pack.Payload))
				return sdkmath.ZeroInt(), nil, sdkerrors.Wrapf(types.ErrInvalidPackage, "payload without header")
			}

			sendSeq, ibcErr := k.CrossChainKeeper.CreateRawIBCPackageWithFee(ctx, pack.ChannelId,
				sdk.FailAckCrossChainPackageType, pack.Payload[sdk.SynPackageHeaderLength:], packageHeader.AckRelayerFee, sdk.NilAckRelayerFee)
			if ibcErr != nil {
				logger.Error("failed to write FailAckCrossChainPackage", "err", err)
				return sdkmath.ZeroInt(), nil, ibcErr
			}
			sendSequence = int64(sendSeq)
		} else if len(result.Payload) != 0 {
			sendSeq, err := k.CrossChainKeeper.CreateRawIBCPackageWithFee(ctx, pack.ChannelId,
				sdk.AckCrossChainPackageType, result.Payload, packageHeader.AckRelayerFee, sdk.NilAckRelayerFee)
			if err != nil {
				logger.Error("failed to write AckCrossChainPackage", "err", err)
				return sdkmath.ZeroInt(), nil, err
			}
			sendSequence = int64(sendSeq)
		}
	}

	claimEvent := &types.EventPackageClaim{
		SrcChainId:      srcChainId,
		DestChainId:     destChainId,
		ChannelId:       uint32(pack.ChannelId),
		PackageType:     uint32(packageHeader.PackageType),
		ReceiveSequence: pack.Sequence,
		SendSequence:    sendSequence,
		RelayerFee:      packageHeader.RelayerFee.String(),
		AckRelayerFee:   packageHeader.AckRelayerFee.String(),
		Crash:           crash,
		ErrorMsg:        result.ErrMsg(),
	}

	return sdkmath.NewIntFromBigInt(packageHeader.RelayerFee), claimEvent, nil
}

func executeClaim(
	ctx sdk.Context,
	app sdk.CrossChainApplication,
	sequence uint64,
	payload []byte,
	header *sdk.PackageHeader,
) (crash bool, result sdk.ExecuteResult) {
	defer func() {
		if r := recover(); r != nil {
			log := fmt.Sprintf("recovered: %v\nstack:\n%v", r, string(debug.Stack()))
			logger := ctx.Logger().With("module", "oracle")
			logger.Error("execute claim panic", "err_log", log)
			crash = true
			result = sdk.ExecuteResult{
				Err: fmt.Errorf("execute claim failed: %v", r),
			}
		}
	}()

	switch header.PackageType {
	case sdk.SynCrossChainPackageType:
		result = app.ExecuteSynPackage(ctx, &sdk.CrossChainAppContext{
			Sequence: sequence,
			Header:   header,
		}, payload[sdk.SynPackageHeaderLength:])
	case sdk.AckCrossChainPackageType:
		result = app.ExecuteAckPackage(ctx, &sdk.CrossChainAppContext{
			Sequence: sequence,
			Header:   header,
		}, payload[sdk.AckPackageHeaderLength:])
	case sdk.FailAckCrossChainPackageType:
		result = app.ExecuteFailAckPackage(ctx, &sdk.CrossChainAppContext{
			Sequence: sequence,
			Header:   header,
		}, payload[sdk.AckPackageHeaderLength:])
	default:
		panic(fmt.Sprintf("receive unexpected package type %d", header.PackageType))
	}
	return crash, result
}
