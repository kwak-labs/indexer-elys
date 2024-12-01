package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerTradeShieldTypes "github.com/elys-network/elys/indexer/txs/tradeshield"
	commonTradeShieldIndxer "github.com/elys-network/elys/indexer/txs/tradeshield/common"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// store params
	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerTradeShieldTypes.MsgUpdateParams{
		Authority: msg.Authority,
		Params: commonTradeShieldIndxer.Params{
			MarketOrderEnabled:   msg.Params.MarketOrderEnabled,
			StakeEnabled:         msg.Params.StakeEnabled,
			ProcessOrdersEnabled: msg.Params.ProcessOrdersEnabled,
			SwapEnabled:          msg.Params.SwapEnabled,
			PerpetualEnabled:     msg.Params.PerpetualEnabled,
			RewardEnabled:        msg.Params.RewardEnabled,
			LeverageEnabled:      msg.Params.LeverageEnabled,
			LimitProcessOrder:    msg.Params.LimitProcessOrder,
			RewardPercentage:     msg.Params.RewardPercentage.String(),
			MarginError:          msg.Params.MarginError.String(),
			MinimumDeposit:       msg.Params.MinimumDeposit.String(),
		},
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateParamsResponse{}, nil
}
