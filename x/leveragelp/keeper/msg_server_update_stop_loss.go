package keeper

import (
	"context"
	"fmt"
	"strconv"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerLeveragelpTypes "github.com/elys-network/elys/indexer/txs/leveragelp"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) UpdateStopLoss(goCtx context.Context, msg *types.MsgUpdateStopLoss) (*types.MsgUpdateStopLossResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	position, found := k.GetPositionWithId(ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Position)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPositionDoesNotExist, fmt.Sprintf("positionId: %d", msg.Position))
	}

	poolId := position.AmmPoolId
	_, found = k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	position.StopLossPrice = msg.Price
	k.SetPosition(ctx, position)

	event := sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(position.Id), 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("collateral", position.Collateral.String()),
		sdk.NewAttribute("liabilities", position.Liabilities.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
		sdk.NewAttribute("stop_loss", position.StopLossPrice.String()),
	)
	ctx.EventManager().EmitEvent(event)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerLeveragelpTypes.MsgUpdateStopLoss{
		Creator:  msg.Creator,
		Position: msg.Position,
		Price:    msg.Price.String(),
		PoolID:   poolId,
		Position_: indexerLeveragelpTypes.PositionStopLoss{
			ID:      position.Id,
			Address: position.Address,
			Collateral: indexerTypes.Token{
				Amount: position.Collateral.Amount.String(),
				Denom:  position.Collateral.Denom,
			},
			Liabilities: position.Liabilities.String(),
			Health:      position.PositionHealth.String(),
			StopLoss:    position.StopLossPrice.String(),
		},
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateStopLossResponse{}, nil
}
