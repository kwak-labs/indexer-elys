package keeper

import (
	"context"
	"strconv"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerLeveragelpTypes "github.com/elys-network/elys/indexer/txs/leveragelp"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) Close(goCtx context.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.Keeper.Close(ctx, msg)
}

func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	closedPosition, repayAmount, err := k.CloseLong(ctx, msg)
	if err != nil {
		return nil, err
	}

	if k.hooks != nil {
		ammPool, err := k.GetAmmPool(ctx, closedPosition.AmmPoolId)
		if err != nil {
			return nil, err
		}
		err = k.hooks.AfterLeverageLpPositionClose(ctx, sdk.MustAccAddressFromBech32(msg.Creator), ammPool)
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(closedPosition.Id), 10)),
		sdk.NewAttribute("address", closedPosition.Address),
		sdk.NewAttribute("collateral", closedPosition.Collateral.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("liabilities", closedPosition.Liabilities.String()),
		sdk.NewAttribute("health", closedPosition.PositionHealth.String()),
	))

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerLeveragelpTypes.MsgClose{
		Creator:     msg.Creator,
		ID:          msg.Id,
		LpAmount:    msg.LpAmount.String(),
		RepayAmount: repayAmount.String(),
		Position: indexerLeveragelpTypes.Position{
			ID:      closedPosition.Id,
			Address: closedPosition.Address,
			Collateral: indexerTypes.Token{
				Amount: closedPosition.Collateral.Amount.String(),
				Denom:  closedPosition.Collateral.Denom,
			},
			Liabilities:    closedPosition.Liabilities.String(),
			PositionHealth: closedPosition.PositionHealth.String(),
		},
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgCloseResponse{}, nil
}
