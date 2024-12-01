package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerLeveragelpTypes "github.com/elys-network/elys/indexer/txs/leveragelp"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	for _, id := range msg.Ids {
		position, err := k.GetPosition(ctx, sender, id)
		if err != nil {
			return nil, err
		}
		err = k.masterchefKeeper.ClaimRewards(ctx, position.GetPositionAddress(), []uint64{position.AmmPoolId}, sender)
		if err != nil {
			return nil, err
		}
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerLeveragelpTypes.MsgClaimRewards{
		Sender: msg.Sender,
		Ids:    msg.Ids,
	}, []string{sender.String()})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgClaimRewardsResponse{}, nil
}
