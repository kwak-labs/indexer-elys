package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) Close(goCtx context.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	// ! Close is handled in the function below
	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/

	return k.Keeper.Close(ctx, msg)
}
