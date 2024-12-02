package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerLeveragelpTypes "github.com/elys-network/elys/indexer/txs/leveragelp"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) Dewhitelist(goCtx context.Context, msg *types.MsgDewhitelist) (*types.MsgDewhitelistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	whitelistAddress := sdk.MustAccAddressFromBech32(msg.WhitelistedAddress)
	k.Keeper.DewhitelistAddress(ctx, whitelistAddress)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerLeveragelpTypes.MsgDewhitelist{
		Authority:          msg.Authority,
		WhitelistedAddress: msg.WhitelistedAddress,
	}, []string{msg.Authority, msg.WhitelistedAddress})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgDewhitelistResponse{}, nil
}
