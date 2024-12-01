package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	// ! Close is handled in the function below
	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	return k.Keeper.Open(ctx, msg)
}
