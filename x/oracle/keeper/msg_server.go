package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"runtime/debug"

	sdkerrors "cosmossdk.io/errors"
	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crosschaintypes "github.com/cosmos/cosmos-sdk/x/crosschain/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
)

type msgServer struct {
	oracleKeeper Keeper
}

// NewMsgServerImpl returns an implementation of the oracle MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{
		oracleKeeper: k,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Claim(goCtx context.Context, req *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	logger := k.oracleKeeper.Logger(ctx)

	// check dest chain id
	if sdk.ChainID(req.DestChainId) != k.oracleKeeper.CrossChainKeeper.GetSrcChainID() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDestChainId, fmt.Sprintf("dest chain id(%d) should be %d", req.DestChainId, k.oracleKeeper.CrossChainKeeper.GetSrcChainID()))
	}

	// check src chain id
	if !k.oracleKeeper.CrossChainKeeper.IsDestChainSupported(sdk.ChainID(req.SrcChainId)) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidSrcChainId, fmt.Sprintf("src chain id(%d) is not supported", req.SrcChainId))
	}

	sequence := k.oracleKeeper.CrossChainKeeper.GetReceiveSequence(ctx, types.RelayPackagesChannelId)
	if sequence != req.Sequence {
		return nil, sdkerrors.Wrapf(types.ErrInvalidReceiveSequence, fmt.Sprintf("current sequence of channel %d is %d", types.RelayPackagesChannelId, sequence))
	}

	signedRelayers, err := k.oracleKeeper.CheckClaim(ctx, req)
	if err != nil {
		return nil, err
	}

	packages := types.Packages{}
	err = rlp.DecodeBytes(req.Payload, &packages)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPayload, "decode payload error")
	}

	events := make([]proto.Message, 0, len(packages))
	totalRelayerFee := big.NewInt(0)
	for idx := range packages {
		pack := packages[idx]

		relayerFee, event, err := handlePackage(ctx, req, k.oracleKeeper, &pack)
		if err != nil {
			// only do log, but let rest package get chance to execute.
			logger.Error("process package failed", "channel", pack.ChannelId, "sequence", pack.Sequence, "error", err.Error())
			return nil, err
		}
		logger.Info("process package success", "channel", pack.ChannelId, "sequence", pack.Sequence)

		events = append(events, event)

		totalRelayerFee = totalRelayerFee.Add(totalRelayerFee, relayerFee)

		// increase channel sequence
		k.oracleKeeper.CrossChainKeeper.IncrReceiveSequence(ctx, pack.ChannelId)
	}

	err = distributeReward(ctx, k.oracleKeeper, req.FromAddress, signedRelayers, totalRelayerFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvents(events...)

	return &types.MsgClaimResponse{}, nil
}

// distributeReward will distribute reward to relayers
func distributeReward(ctx sdk.Context, oracleKeeper Keeper, relayer string, signedRelayers []string, relayerFee *big.Int) error {
	if relayerFee.Cmp(big.NewInt(0)) <= 0 {
		oracleKeeper.Logger(ctx).Info("total relayer fee is zero")
		return nil
	}

	relayerAddr, err := sdk.AccAddressFromHexUnsafe(relayer)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidAddress, fmt.Sprintf("relayer address (%s) is invalid", relayer))
	}

	otherRelayers := make([]sdk.AccAddress, 0, len(signedRelayers))
	for idx := range signedRelayers {
		signedRelayerAddr, err := sdk.AccAddressFromHexUnsafe(signedRelayers[idx])
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidAddress, fmt.Sprintf("relayer address (%s) is invalid", relayer))
		}
		if !signedRelayerAddr.Equals(relayerAddr) {
			otherRelayers = append(otherRelayers, relayerAddr)
		}
	}

	totalDistributed, otherRelayerReward := big.NewInt(0), big.NewInt(0)

	relayerRewardShare := oracleKeeper.GetRelayerRewardShare(ctx)

	// calculate the reward to distribute to each other relayer
	if len(otherRelayers) > 0 {
		otherRelayerReward = otherRelayerReward.Mul(big.NewInt(100-int64(relayerRewardShare)), relayerFee)
		otherRelayerReward = otherRelayerReward.Div(otherRelayerReward, big.NewInt(100))
		otherRelayerReward = otherRelayerReward.Div(otherRelayerReward, big.NewInt(int64(len(otherRelayers))))
	}

	bondDenom := oracleKeeper.StakingKeeper.BondDenom(ctx)
	if otherRelayerReward.Cmp(big.NewInt(0)) > 0 {
		for idx := range otherRelayers {
			err = oracleKeeper.BankKeeper.SendCoinsFromModuleToAccount(ctx,
				crosschaintypes.ModuleName,
				otherRelayers[idx],
				sdk.Coins{sdk.Coin{Denom: bondDenom, Amount: sdk.NewIntFromBigInt(otherRelayerReward)}},
			)
			if err != nil {
				return err
			}
			totalDistributed = totalDistributed.Add(totalDistributed, otherRelayerReward)
		}
	}

	remainingReward := relayerFee.Sub(relayerFee, totalDistributed)
	if remainingReward.Cmp(big.NewInt(0)) > 0 {
		err = oracleKeeper.BankKeeper.SendCoinsFromModuleToAccount(ctx,
			crosschaintypes.ModuleName,
			relayerAddr,
			sdk.Coins{sdk.Coin{Denom: bondDenom, Amount: sdk.NewIntFromBigInt(remainingReward)}},
		)
		if err != nil {
			return err
		}
	}

	return err
}

func handlePackage(
	ctx sdk.Context,
	req *types.MsgClaim,
	oracleKeeper Keeper,
	pack *types.Package,
) (*big.Int, *types.EventPackageClaim, error) {
	logger := oracleKeeper.Logger(ctx)

	crossChainApp := oracleKeeper.CrossChainKeeper.GetCrossChainApp(pack.ChannelId)
	if crossChainApp == nil {
		return nil, nil, sdkerrors.Wrapf(types.ErrChannelNotRegistered, "channel %d not registered", pack.ChannelId)
	}

	sequence := oracleKeeper.CrossChainKeeper.GetReceiveSequence(ctx, pack.ChannelId)
	if sequence != pack.Sequence {
		return nil, nil, sdkerrors.Wrapf(types.ErrInvalidReceiveSequence,
			fmt.Sprintf("current sequence of channel %d is %d", pack.ChannelId, sequence))
	}

	packageHeader, err := sdk.DecodePackageHeader(pack.Payload)
	if err != nil {
		return nil, nil, sdkerrors.Wrapf(types.ErrInvalidPayloadHeader, "payload header is invalid")
	}

	if packageHeader.Timestamp != req.Timestamp {
		return nil, nil, sdkerrors.Wrapf(types.ErrInvalidPayloadHeader,
			"timestamp(%d) is not the same in payload header(%d)", req.Timestamp, packageHeader.Timestamp)
	}

	if !sdk.IsValidCrossChainPackageType(packageHeader.PackageType) {
		return nil, nil, sdkerrors.Wrapf(types.ErrInvalidPackageType,
			fmt.Sprintf("package type %d is invalid", packageHeader.PackageType))
	}

	cacheCtx, write := ctx.CacheContext()
	crash, result := executeClaim(cacheCtx, crossChainApp, pack.Payload, packageHeader.PackageType, packageHeader.SynRelayerFee)
	if result.IsOk() {
		write()
	}

	// write ack package
	var sendSequence int64 = -1
	if packageHeader.PackageType == sdk.SynCrossChainPackageType {
		if crash {
			var ibcErr error
			var sendSeq uint64
			if len(pack.Payload) >= sdk.SynPackageHeaderLength {
				sendSeq, ibcErr = oracleKeeper.CrossChainKeeper.CreateRawIBCPackageWithFee(ctx, pack.ChannelId,
					sdk.FailAckCrossChainPackageType, pack.Payload[sdk.SynPackageHeaderLength:], packageHeader.AckRelayerFee, sdk.NilAckRelayerFee)
			} else {
				logger.Error("found payload without header",
					"channelID", pack.ChannelId, "sequence", pack.Sequence, "payload", hex.EncodeToString(pack.Payload))
				return nil, nil, sdkerrors.Wrapf(types.ErrInvalidPackage, "payload without header")
			}

			if ibcErr != nil {
				logger.Error("failed to write FailAckCrossChainPackage", "err", err)
				return nil, nil, ibcErr
			}
			sendSequence = int64(sendSeq)
		} else if len(result.Payload) != 0 {
			sendSeq, err := oracleKeeper.CrossChainKeeper.CreateRawIBCPackageWithFee(ctx, pack.ChannelId,
				sdk.AckCrossChainPackageType, result.Payload, packageHeader.AckRelayerFee, sdk.NilAckRelayerFee)
			if err != nil {
				logger.Error("failed to write AckCrossChainPackage", "err", err)
				return nil, nil, err
			}
			sendSequence = int64(sendSeq)
		}
	}

	claimEvent := &types.EventPackageClaim{
		SrcChainId:      req.SrcChainId,
		DestChainId:     req.DestChainId,
		ChannelId:       uint32(pack.ChannelId),
		PackageType:     uint32(packageHeader.PackageType),
		ReceiveSequence: pack.Sequence,
		SendSequence:    sendSequence,
		SynRelayFee:     packageHeader.SynRelayerFee.String(),
		AckRelayFee:     packageHeader.AckRelayerFee.String(),
		Crash:           crash,
		ErrorMsg:        result.ErrMsg(),
	}

	return packageHeader.SynRelayerFee, claimEvent, nil
}

func executeClaim(
	ctx sdk.Context,
	app sdk.CrossChainApplication,
	payload []byte,
	packageType sdk.CrossChainPackageType,
	relayerFee *big.Int,
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

	switch packageType {
	case sdk.SynCrossChainPackageType:
		result = app.ExecuteSynPackage(ctx, payload[sdk.SynPackageHeaderLength:], relayerFee)
	case sdk.AckCrossChainPackageType:
		result = app.ExecuteAckPackage(ctx, payload[sdk.AckPackageHeaderLength:])
	case sdk.FailAckCrossChainPackageType:
		result = app.ExecuteFailAckPackage(ctx, payload[sdk.AckPackageHeaderLength:])
	default:
		panic(fmt.Sprintf("receive unexpected package type %d", packageType))
	}
	return
}
