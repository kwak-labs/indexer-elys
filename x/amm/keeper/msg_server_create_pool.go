package keeper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerAmmTypes "github.com/elys-network/elys/indexer/txs/amm"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)

	if !params.IsCreatorAllowed(msg.Sender) {
		return nil, fmt.Errorf("sender is not allowed to create pool")
	}

	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	baseAssetExists := false
	for _, asset := range msg.PoolAssets {
		baseAssetExists = k.CheckBaseAssetExist(ctx, asset.Token.Denom)
		if baseAssetExists {
			break
		}
	}
	if !baseAssetExists {
		return nil, errorsmod.Wrapf(types.ErrOnlyBaseAssetsPoolAllowed, "one of the asset must be from %s", strings.Join(params.BaseAssets, ", "))
	}

	feeAssetExists := k.CheckBaseAssetExist(ctx, msg.PoolParams.FeeDenom)
	if !feeAssetExists {
		return nil, fmt.Errorf("fee denom must be from %s", strings.Join(params.BaseAssets, ", "))
	}

	if !params.PoolCreationFee.IsNil() && params.PoolCreationFee.IsPositive() {
		feeCoins := sdk.Coins{sdk.NewCoin(ptypes.Elys, params.PoolCreationFee)}
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, feeCoins); err != nil {
			return nil, err
		}
	}

	poolId, err := k.Keeper.CreatePool(ctx, msg)
	if err != nil {
		return &types.MsgCreatePoolResponse{}, err
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	poolAssets := make([]indexerAmmTypes.PoolAsset, len(msg.PoolAssets))
	for i, asset := range msg.PoolAssets {
		poolAssets[i] = indexerAmmTypes.PoolAsset{
			Token: indexerTypes.Token{
				Amount: asset.Token.Amount.String(),
				Denom:  asset.Token.Denom,
			},
		}
	}

	indexer.QueueTransaction(ctx, indexerAmmTypes.MsgCreatePool{
		Sender: msg.Sender,
		PoolParams: indexerAmmTypes.PoolParams{
			SwapFee:   msg.PoolParams.SwapFee.String(),
			UseOracle: msg.PoolParams.UseOracle,
			FeeDenom:  msg.PoolParams.FeeDenom,
		},
		PoolAssets: poolAssets,
		PoolID:     poolId,
	}, []string{})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtPoolCreated,
			sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgCreatePoolResponse{
		PoolID: poolId,
	}, nil
}
