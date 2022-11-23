package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"runtime/debug"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
	"github.com/gogo/protobuf/proto"
)

type msgServer struct {
	oracleKeeper  Keeper
	stakingKeeper types.StakingKeeper
}

// NewMsgServerImpl returns an implementation of the oracle MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper, stakingKeeper types.StakingKeeper) types.MsgServer {
	return &msgServer{
		oracleKeeper:  k,
		stakingKeeper: stakingKeeper,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Claim(goCtx context.Context, req *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sequence := k.oracleKeeper.CrossChainKeeper.GetReceiveSequence(ctx, sdk.ChainID(req.ChainId), types.RelayPackagesChannelId)
	if sequence != req.Sequence {
		return nil, sdkerrors.Wrapf(types.ErrInvalidReceiveSequence, fmt.Sprintf("current sequence of channel %d is %d", types.RelayPackagesChannelId, sequence))
	}

	err := k.oracleKeeper.ProcessClaim(ctx, req)
	if err != nil {
		return nil, err
	}

	packages := types.Packages{}
	err = rlp.DecodeBytes(req.Payload, &packages)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPayload, fmt.Sprintf("decode payload error"))
	}

	events := make([]proto.Message, 0, len(packages))
	for _, pack := range packages {
		event, err := handlePackage(ctx, k.oracleKeeper, sdk.ChainID(req.ChainId), &pack)
		if err != nil {
			// only do log, but let reset package get chance to execute.
			ctx.Logger().With("module", "oracle").Error(fmt.Sprintf("process package failed, channel=%d, sequence=%d, error=%v", pack.ChannelId, pack.Sequence, err))
			return nil, err
		}
		ctx.Logger().With("module", "oracle").Info(fmt.Sprintf("process package success, channel=%d, sequence=%d", pack.ChannelId, pack.Sequence))

		events = append(events, event)

		// increase channel sequence
		k.oracleKeeper.CrossChainKeeper.IncrReceiveSequence(ctx, sdk.ChainID(req.ChainId), pack.ChannelId)
	}

	ctx.EventManager().EmitTypedEvents(events...)

	return &types.MsgClaimResponse{}, nil
}

func handlePackage(ctx sdk.Context, oracleKeeper Keeper, chainId sdk.ChainID, pack *types.Package) (*types.EventPackageClaim, error) {
	logger := ctx.Logger().With("module", "x/oracle")

	crossChainApp := oracleKeeper.CrossChainKeeper.GetCrossChainApp(pack.ChannelId)
	if crossChainApp == nil {
		return nil, sdkerrors.Wrapf(types.ErrChannelNotRegistered, "channel %d not registered", pack.ChannelId)
	}

	sequence := oracleKeeper.CrossChainKeeper.GetReceiveSequence(ctx, chainId, pack.ChannelId)
	if sequence != pack.Sequence {
		return nil, sdkerrors.Wrapf(types.ErrInvalidReceiveSequence, fmt.Sprintf("current sequence of channel %d is %d", pack.ChannelId, sequence))
	}

	packageType, _, relayFee, err := sdk.DecodePackageHeader(pack.Payload)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPayloadHeader, "payload header is invalid")
	}

	if !sdk.IsValidCrossChainPackageType(packageType) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPackageType, fmt.Sprintf("pcakge type %d is invalid", packageType))
	}

	feeAmount := relayFee.Int64()
	if feeAmount < 0 {
		return nil, sdkerrors.Wrapf(types.ErrFeeOverflow, fmt.Sprintf("fee(%s) is overflow", relayFee.String()))
	}

	fee := sdk.Coins{sdk.Coin{Denom: sdk.NativeTokenSymbol, Amount: sdk.NewInt(feeAmount)}}
	err = oracleKeeper.SendCoinsFromAccountToFeeCollector(ctx, sdk.PegAccount, fee)
	if err != nil {
		return nil, err
	}

	cacheCtx, write := ctx.CacheContext()
	crash, result := executeClaim(cacheCtx, crossChainApp, pack.Payload, packageType, feeAmount)
	if result.IsOk() {
		write()
	} else {
		oracleKeeper.Metrics.ErrNumOfChannels.With("channel_id", fmt.Sprintf("%d", pack.ChannelId)).Add(1)
	}

	// write ack package
	var sendSequence int64 = -1
	if packageType == sdk.SynCrossChainPackageType {
		if crash {
			var ibcErr error
			var sendSeq uint64
			if len(pack.Payload) >= sdk.PackageHeaderLength {
				sendSeq, ibcErr = oracleKeeper.CrossChainKeeper.CreateRawIBCPackage(ctx, chainId,
					pack.ChannelId, sdk.FailAckCrossChainPackageType, pack.Payload[sdk.PackageHeaderLength:])
			} else {
				logger.Error("found payload without header", "channelID", pack.ChannelId, "sequence", pack.Sequence, "payload", hex.EncodeToString(pack.Payload))
				return nil, sdkerrors.Wrapf(types.ErrInvalidPackage, fmt.Sprintf("payload without header"))
			}

			if ibcErr != nil {
				logger.Error("failed to write FailAckCrossChainPackage", "err", err)
				return nil, ibcErr
			}
			sendSequence = int64(sendSeq)
		} else {
			if len(result.Payload) != 0 {
				sendSeq, err := oracleKeeper.CrossChainKeeper.CreateRawIBCPackage(ctx, chainId,
					pack.ChannelId, sdk.AckCrossChainPackageType, result.Payload)
				if err != nil {
					logger.Error("failed to write AckCrossChainPackage", "err", err)
					return nil, err
				}
				sendSequence = int64(sendSeq)
			}
		}
	}

	claimEvent := &types.EventPackageClaim{
		ChainId:         uint32(chainId),
		ChannelId:       uint32(pack.ChannelId),
		PackageType:     uint32(packageType),
		ReceiveSequence: pack.Sequence,
		SendSequence:    sendSequence,
		Fee:             feeAmount,
		Crash:           crash,
		ErrorMsg:        result.ErrMsg(),
	}

	return claimEvent, nil
}

func executeClaim(ctx sdk.Context, app sdk.CrossChainApplication, payload []byte, packageType sdk.CrossChainPackageType, relayerFee int64) (crash bool, result sdk.ExecuteResult) {
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
		result = app.ExecuteSynPackage(ctx, payload[sdk.PackageHeaderLength:], relayerFee)
	case sdk.AckCrossChainPackageType:
		result = app.ExecuteAckPackage(ctx, payload[sdk.PackageHeaderLength:])
	case sdk.FailAckCrossChainPackageType:
		result = app.ExecuteFailAckPackage(ctx, payload[sdk.PackageHeaderLength:])
	default:
		panic(fmt.Sprintf("receive unexpected package type %d", packageType))
	}
	return
}
