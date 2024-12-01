package keeper

import (
	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerOracleTypes "github.com/elys-network/elys/indexer/txs/oracle"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) FeedMultiplePrices(goCtx context.Context, msg *types.MsgFeedMultiplePrices) (*types.MsgFeedMultiplePricesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	feeder, found := k.Keeper.GetPriceFeeder(ctx, creator)
	if !found {
		return nil, types.ErrNotAPriceFeeder
	}

	if !feeder.IsActive {
		return nil, types.ErrPriceFeederNotActive
	}

	for _, feedPrice := range msg.FeedPrices {
		price := types.Price{
			Asset:       feedPrice.Asset,
			Price:       feedPrice.Price,
			Source:      feedPrice.Source,
			Provider:    msg.Creator,
			Timestamp:   uint64(ctx.BlockTime().Unix()),
			BlockHeight: uint64(ctx.BlockHeight()),
		}
		k.SetPrice(ctx, price)
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	// Convert feed prices to indexer format
	indexerFeedPrices := make([]indexerOracleTypes.FeedPrice, len(msg.FeedPrices))
	for i, fp := range msg.FeedPrices {
		indexerFeedPrices[i] = indexerOracleTypes.FeedPrice{
			Asset:  fp.Asset,
			Price:  fp.Price.String(),
			Source: fp.Source,
		}
	}

	// Queue the transaction
	indexer.QueueTransaction(ctx, indexerOracleTypes.MsgFeedMultiplePrices{
		Creator:    msg.Creator,
		FeedPrices: indexerFeedPrices,
		Timestamp:  uint64(ctx.BlockTime().Unix()),
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgFeedMultiplePricesResponse{}, nil
}
