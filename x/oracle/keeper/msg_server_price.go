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

func (k msgServer) FeedPrice(goCtx context.Context, msg *types.MsgFeedPrice) (*types.MsgFeedPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	provider := sdk.MustAccAddressFromBech32(msg.Provider)
	feeder, found := k.Keeper.GetPriceFeeder(ctx, provider)
	if !found {
		return nil, types.ErrNotAPriceFeeder
	}

	if !feeder.IsActive {
		return nil, types.ErrPriceFeederNotActive
	}

	price := types.Price{
		Asset:       msg.FeedPrice.Asset,
		Price:       msg.FeedPrice.Price,
		Source:      msg.FeedPrice.Source,
		Provider:    msg.Provider,
		Timestamp:   uint64(ctx.BlockTime().Unix()),
		BlockHeight: uint64(ctx.BlockHeight()),
	}

	k.SetPrice(ctx, price)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerOracleTypes.MsgFeedPrice{
		Provider:    msg.Provider,
		Asset:       msg.FeedPrice.Asset,
		Price:       msg.FeedPrice.Price.String(),
		Source:      msg.FeedPrice.Source,
		Timestamp:   uint64(ctx.BlockTime().Unix()),
		BlockHeight: uint64(ctx.BlockHeight()),
	}, []string{msg.Provider})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgFeedPriceResponse{}, nil
}
