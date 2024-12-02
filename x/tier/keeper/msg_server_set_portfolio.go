package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerTierTypes "github.com/elys-network/elys/indexer/txs/tier"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tier/types"
)

func (k msgServer) SetPortfolio(goCtx context.Context, msg *types.MsgSetPortfolio) (*types.MsgSetPortfolioResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	user := sdk.MustAccAddressFromBech32(msg.User)
	k.RetrieveAllPortfolio(ctx, user)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerTierTypes.MsgSetPortfolio{
		Creator: msg.Creator,
		User:    msg.User,
	}, []string{msg.Creator, msg.User})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgSetPortfolioResponse{}, nil
}
