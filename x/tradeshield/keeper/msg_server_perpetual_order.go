package keeper

import (
	"context"
	"fmt"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerTradeshieldTypes "github.com/elys-network/elys/indexer/txs/tradeshield"
	"github.com/elys-network/elys/indexer/txs/tradeshield/common"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) CreatePerpetualOpenOrder(goCtx context.Context, msg *types.MsgCreatePerpetualOpenOrder) (*types.MsgCreatePerpetualOpenOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify if perpetual pool exists
	_, found := k.perpetual.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("pool %d not found", msg.PoolId))
	}

	var pendingPerpetualOrder = types.PerpetualOrder{
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		TriggerPrice:       msg.TriggerPrice,
		Collateral:         msg.Collateral,
		OwnerAddress:       msg.OwnerAddress,
		TradingAsset:       msg.TradingAsset,
		Position:           msg.Position,
		Leverage:           msg.Leverage,
		TakeProfitPrice:    msg.TakeProfitPrice,
		StopLossPrice:      msg.StopLossPrice,
		PoolId:             msg.PoolId,
		PositionId:         0,
		Status:             types.Status_PENDING,
	}

	// Verify if user hasn't created a order for same pool with pending status
	// Note: A user can have either
	// at most one pending order for a pool
	// or a position in the pool
	pendingStatus := types.Status_PENDING
	orders, _, err := k.GetPendingPerpetualOrdersForAddress(ctx, msg.OwnerAddress, &pendingStatus, nil)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		if order.PoolId == msg.PoolId && order.Position == msg.Position &&
			order.Collateral.Denom == msg.Collateral.Denom &&
			order.TradingAsset == msg.TradingAsset {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "user already has a order for the same pool")
		}
	}

	// Verify if user doesn't have a position in the same pool
	// Should not create a order for a position where the user already has a position in the same pool
	mtps, _, err := k.perpetual.GetMTPsForAddressWithPagination(ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), nil)
	if err != nil {
		return nil, err
	}
	for _, mtp := range mtps {
		if mtp.Mtp.AmmPoolId == msg.PoolId && mtp.Mtp.Position == perpetualtypes.Position(msg.Position) && mtp.Mtp.CollateralAsset == msg.Collateral.Denom && mtp.Mtp.TradingAsset == msg.TradingAsset {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "user already has a position in the same pool")
		}
	}

	id := k.AppendPendingPerpetualOrder(
		ctx,
		pendingPerpetualOrder,
	)

	// Verify if order is valid before saving
	_, err = k.perpetual.HandleOpenEstimation(ctx, &perpetualtypes.QueryOpenEstimationRequest{
		Position:        perpetualtypes.Position(msg.Position),
		Leverage:        msg.Leverage,
		TradingAsset:    msg.TradingAsset,
		Collateral:      msg.Collateral,
		TakeProfitPrice: msg.TakeProfitPrice,
		PoolId:          msg.PoolId,
		LimitPrice:      msg.TriggerPrice.Rate,
	})
	if err != nil {
		return nil, err
	}

	// set the order id
	pendingPerpetualOrder.OrderId = id

	// send collateral amount from owner to the order address
	ownerAddress := sdk.MustAccAddressFromBech32(pendingPerpetualOrder.OwnerAddress)
	err = k.Keeper.bank.SendCoins(ctx, ownerAddress, pendingPerpetualOrder.GetOrderAddress(), sdk.NewCoins(pendingPerpetualOrder.Collateral))
	if err != nil {
		return nil, err
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerTradeshieldTypes.MsgCreatePerpetualOpenOrder{
		OwnerAddress: msg.OwnerAddress,
		TriggerPrice: common.TriggerPrice{
			TradingAssetDenom: msg.TriggerPrice.TradingAssetDenom,
			Rate:              msg.TriggerPrice.Rate.String(),
		},
		Collateral: indexerTypes.Token{
			Amount: msg.Collateral.Amount.String(),
			Denom:  msg.Collateral.Denom,
		},
		TradingAsset:    msg.TradingAsset,
		Position:        int32(msg.Position),
		Leverage:        msg.Leverage.String(),
		TakeProfitPrice: msg.TakeProfitPrice.String(),
		StopLossPrice:   msg.StopLossPrice.String(),
		PoolID:          msg.PoolId,
		OrderID:         id,
	}, []string{msg.OwnerAddress})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgCreatePerpetualOpenOrderResponse{
		OrderId: pendingPerpetualOrder.OrderId,
	}, nil
}

func (k msgServer) CreatePerpetualCloseOrder(goCtx context.Context, msg *types.MsgCreatePerpetualCloseOrder) (*types.MsgCreatePerpetualCloseOrderResponse, error) {
	// Disable for v1
	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "disabled for v1")
	// ctx := sdk.UnwrapSDKContext(goCtx)

	// // check if the position owner address matches the msg owner address
	// position, err := k.perpetual.GetMTP(ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), msg.PositionId)
	// if err != nil {
	// 	return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("position %d not found", msg.PositionId))
	// }

	// var pendingPerpetualOrder = types.PerpetualOrder{
	// 	PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
	// 	TriggerPrice: types.TriggerPrice{
	// 		TradingAssetDenom: position.TradingAsset,
	// 		Rate:              msg.TriggerPrice.Rate,
	// 	},
	// 	OwnerAddress: position.Address,
	// 	PositionId:   position.Id,
	// }

	// id := k.AppendPendingPerpetualOrder(
	// 	ctx,
	// 	pendingPerpetualOrder,
	// )

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	// indexer.QueueTransaction(ctx, indexerTradeshieldTypes.MsgCreatePerpetualCloseOrder{
	// 	OwnerAddress: msg.OwnerAddress,
	// 	TriggerPrice: indexerTradeshieldTypes.TriggerPrice{
	// 		TradingAssetDenom: msg.TriggerPrice.TradingAssetDenom,
	// 		Rate:             msg.TriggerPrice.Rate.String(),
	// 	},
	// 	PositionID: msg.PositionId,
	// 	OrderID:    id,
	// }, []string{msg.OwnerAddress})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	// return &types.MsgCreatePerpetualCloseOrderResponse{
	// 	OrderId: id,
	// }, nil
}

func (k msgServer) UpdatePerpetualOrder(goCtx context.Context, msg *types.MsgUpdatePerpetualOrder) (*types.MsgUpdatePerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	order, found := k.GetPendingPerpetualOrder(ctx, msg.OrderId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.OrderId))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.OwnerAddress != order.OwnerAddress {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	perpetualParams := k.perpetual.GetParams(ctx)

	ratio := order.TakeProfitPrice.Quo(msg.TriggerPrice.Rate)
	if order.Position == types.PerpetualPosition_LONG {
		if ratio.LT(perpetualParams.MinimumLongTakeProfitPriceRatio) || ratio.GT(perpetualParams.MaximumLongTakeProfitPriceRatio) {
			return nil, fmt.Errorf("invalid trigger price, take profit price should be between %s and %s times of current market price for long (current ratio: %s)", perpetualParams.MinimumLongTakeProfitPriceRatio.String(), perpetualParams.MaximumLongTakeProfitPriceRatio.String(), ratio.String())
		}
	}
	if order.Position == types.PerpetualPosition_SHORT {
		if ratio.GT(perpetualParams.MaximumShortTakeProfitPriceRatio) {
			return nil, fmt.Errorf("invalid trigger price, take profit price should be less than %s times of current market price for short (current ratio: %s)", perpetualParams.MaximumShortTakeProfitPriceRatio.String(), ratio.String())
		}
	}

	// update the order
	order.TriggerPrice = msg.TriggerPrice
	k.SetPendingPerpetualOrder(ctx, order)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerTradeshieldTypes.MsgUpdatePerpetualOrder{
		OwnerAddress: msg.OwnerAddress,
		OrderID:      msg.OrderId,
		TriggerPrice: common.TriggerPrice{
			TradingAssetDenom: msg.TriggerPrice.TradingAssetDenom,
			Rate:              msg.TriggerPrice.Rate.String(),
		},
	}, []string{msg.OwnerAddress})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdatePerpetualOrderResponse{}, nil
}

func (k msgServer) CancelPerpetualOrder(goCtx context.Context, msg *types.MsgCancelPerpetualOrder) (*types.MsgCancelPerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	order, found := k.GetPendingPerpetualOrder(ctx, msg.OrderId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("order %d doesn't exist", msg.OrderId))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.OwnerAddress != order.OwnerAddress {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// send the collateral amount back to the owner
	ownerAddress := sdk.MustAccAddressFromBech32(order.OwnerAddress)
	err := k.Keeper.bank.SendCoins(ctx, order.GetOrderAddress(), ownerAddress, sdk.NewCoins(order.Collateral))
	if err != nil {
		return nil, err
	}

	k.RemovePendingPerpetualOrder(ctx, msg.OrderId)
	types.EmitCancelPerpetualOrderEvent(ctx, order)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerTradeshieldTypes.MsgCancelPerpetualOrder{
		OwnerAddress: msg.OwnerAddress,
		OrderID:      msg.OrderId,
		Collateral: indexerTypes.Token{
			Amount: order.Collateral.Amount.String(),
			Denom:  order.Collateral.Denom,
		},
	}, []string{msg.OwnerAddress})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgCancelPerpetualOrderResponse{
		OrderId: order.OrderId,
	}, nil
}

func (k msgServer) CancelPerpetualOrders(goCtx context.Context, msg *types.MsgCancelPerpetualOrders) (*types.MsgCancelPerpetualOrdersResponse, error) {
	if len(msg.OrderIds) == 0 {
		return nil, types.ErrSizeZero
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	ctx := sdk.UnwrapSDKContext(goCtx)
	indexer.QueueTransaction(ctx, indexerTradeshieldTypes.MsgCancelPerpetualOrders{
		OwnerAddress: msg.OwnerAddress,
		OrderIds:     msg.OrderIds,
	}, []string{msg.OwnerAddress})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	// loop through the spot orders and cancel them
	for _, orderId := range msg.OrderIds {
		_, err := k.CancelPerpetualOrder(goCtx, &types.MsgCancelPerpetualOrder{
			OwnerAddress: msg.OwnerAddress,
			OrderId:      orderId,
		})
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgCancelPerpetualOrdersResponse{}, nil
}
