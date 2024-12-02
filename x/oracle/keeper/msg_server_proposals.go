package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerOracleTypes "github.com/elys-network/elys/indexer/txs/oracle"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.Keeper.SetParams(ctx, msg.Params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	feeLimit := make([]indexerTypes.Token, len(msg.Params.FeeLimit))
	for i, coin := range msg.Params.FeeLimit {
		feeLimit[i] = indexerTypes.Token{
			Amount: coin.Amount.String(),
			Denom:  coin.Denom,
		}
	}

	indexer.QueueTransaction(ctx, indexerOracleTypes.MsgUpdateParams{
		Authority: msg.Authority,
		Params: indexerOracleTypes.Params{
			BandChannelSource: msg.Params.BandChannelSource,
			OracleScriptID:    msg.Params.OracleScriptID,
			Multiplier:        msg.Params.Multiplier,
			AskCount:          msg.Params.AskCount,
			MinCount:          msg.Params.MinCount,
			FeeLimit:          feeLimit,
			PrepareGas:        msg.Params.PrepareGas,
			ExecuteGas:        msg.Params.ExecuteGas,
			ClientID:          msg.Params.ClientID,
			BandEpoch:         msg.Params.BandEpoch,
			PriceExpiryTime:   msg.Params.PriceExpiryTime,
			LifeTimeInBlocks:  msg.Params.LifeTimeInBlocks,
		},
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) RemoveAssetInfo(goCtx context.Context, msg *types.MsgRemoveAssetInfo) (*types.MsgRemoveAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.Keeper.RemoveAssetInfo(ctx, msg.Denom)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerOracleTypes.MsgRemoveAssetInfo{
		Authority: msg.Authority,
		Denom:     msg.Denom,
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgRemoveAssetInfoResponse{}, nil
}

func (k msgServer) AddPriceFeeders(goCtx context.Context, msg *types.MsgAddPriceFeeders) (*types.MsgAddPriceFeedersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	for _, feeder := range msg.Feeders {
		k.Keeper.SetPriceFeeder(ctx, types.PriceFeeder{
			Feeder:   feeder,
			IsActive: true,
		})
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerOracleTypes.MsgAddPriceFeeders{
		Authority: msg.Authority,
		Feeders:   msg.Feeders,
	}, append([]string{msg.Authority}, msg.Feeders...))
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgAddPriceFeedersResponse{}, nil
}

func (k msgServer) RemovePriceFeeders(goCtx context.Context, msg *types.MsgRemovePriceFeeders) (*types.MsgRemovePriceFeedersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	for _, feeder := range msg.Feeders {
		k.Keeper.RemovePriceFeeder(ctx, sdk.MustAccAddressFromBech32(feeder))
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerOracleTypes.MsgRemovePriceFeeders{
		Authority: msg.Authority,
		Feeders:   msg.Feeders,
	}, append([]string{msg.Authority}, msg.Feeders...))
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgRemovePriceFeedersResponse{}, nil
}
