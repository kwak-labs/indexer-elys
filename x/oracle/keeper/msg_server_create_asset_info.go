package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerOracleTypes "github.com/elys-network/elys/indexer/txs/oracle"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) CreateAssetInfo(goCtx context.Context, msg *types.MsgCreateAssetInfo) (*types.MsgCreateAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, found := k.GetAssetInfo(ctx, msg.Denom)

	if found {
		return nil, errors.Wrapf(types.ErrAssetWasCreated, "%s", msg.Denom)
	}

	k.Keeper.SetAssetInfo(ctx, types.AssetInfo{
		Denom:      msg.Denom,
		Display:    msg.Display,
		BandTicker: msg.BandTicker,
		ElysTicker: msg.ElysTicker,
		Decimal:    msg.Decimal,
	})

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerOracleTypes.MsgCreateAssetInfo{
		Creator:    msg.Creator,
		Denom:      msg.Denom,
		Display:    msg.Display,
		BandTicker: msg.BandTicker,
		ElysTicker: msg.ElysTicker,
		Decimal:    msg.Decimal,
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgCreateAssetInfoResponse{}, nil
}
