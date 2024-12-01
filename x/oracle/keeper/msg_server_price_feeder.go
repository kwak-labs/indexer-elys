package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerOracleTypes "github.com/elys-network/elys/indexer/txs/oracle"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) SetPriceFeeder(goCtx context.Context, msg *types.MsgSetPriceFeeder) (*types.MsgSetPriceFeederResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	feederAccount := sdk.MustAccAddressFromBech32(msg.Feeder)
	_, found := k.Keeper.GetPriceFeeder(ctx, feederAccount)
	if !found {
		return nil, types.ErrNotAPriceFeeder
	}
	k.Keeper.SetPriceFeeder(ctx, types.PriceFeeder{
		Feeder:   msg.Feeder,
		IsActive: msg.IsActive,
	})

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerOracleTypes.MsgSetPriceFeeder{
		Feeder:   msg.Feeder,
		IsActive: msg.IsActive,
	}, []string{msg.Feeder})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgSetPriceFeederResponse{}, nil
}

func (k msgServer) DeletePriceFeeder(goCtx context.Context, msg *types.MsgDeletePriceFeeder) (*types.MsgDeletePriceFeederResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	feederAccount := sdk.MustAccAddressFromBech32(msg.Feeder)
	_, found := k.Keeper.GetPriceFeeder(ctx, feederAccount)
	if !found {
		return nil, types.ErrNotAPriceFeeder
	}

	k.RemovePriceFeeder(ctx, feederAccount)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerOracleTypes.MsgDeletePriceFeeder{
		Feeder: msg.Feeder,
	}, []string{msg.Feeder})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgDeletePriceFeederResponse{}, nil
}
