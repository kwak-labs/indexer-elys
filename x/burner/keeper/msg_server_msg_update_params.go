package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerBurnerTypes "github.com/elys-network/elys/indexer/txs/burner"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/burner/types"
)

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	if len(msg.Params.EpochIdentifier) == 0 {
		return nil, types.ErrInvalidEpochIdentifier
	}

	params := k.GetParams(ctx)
	params.EpochIdentifier = msg.Params.EpochIdentifier
	k.SetParams(ctx, &params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implemsentation*/
	indexer.QueueTransaction(ctx, indexerBurnerTypes.MsgUpdateParams{
		Authority: msg.Authority,
		Params: indexerBurnerTypes.Params{
			EpochIdentifier: msg.Params.EpochIdentifier,
		},
	}, []string{})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateParamsResponse{}, nil
}
