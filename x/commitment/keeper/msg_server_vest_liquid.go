package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerCommitmentsTypes "github.com/elys-network/elys/indexer/txs/commitments"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// VestLiquid converts user's balance to vesting to be utilized for normal tokens vesting like ATOM vesting
func (k msgServer) VestLiquid(goCtx context.Context, msg *types.MsgVestLiquid) (*types.MsgVestLiquidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	creator := sdk.MustAccAddressFromBech32(msg.Creator)

	if err := k.DepositLiquidTokensClaimed(ctx, msg.Denom, msg.Amount, creator); err != nil {
		return &types.MsgVestLiquidResponse{}, err
	}

	if err := k.ProcessTokenVesting(ctx, msg.Denom, msg.Amount, creator); err != nil {
		return &types.MsgVestLiquidResponse{}, err
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerCommitmentsTypes.MsgVestLiquid{
		Creator: creator.String(),
		Token: indexerTypes.Token{
			Amount: msg.Amount.String(),
			Denom:  msg.Denom,
		},
	}, []string{creator.String()})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgVestLiquidResponse{}, nil
}
