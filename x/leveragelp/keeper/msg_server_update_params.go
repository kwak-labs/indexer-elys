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

// Update params through gov proposal
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
	indexer.QueueTransaction(ctx, indexerLeveragelpTypes.MsgUpdateParams{
		Authority: msg.Authority,
		Params: indexerLeveragelpTypes.Params{
			LeverageMax:         msg.Params.LeverageMax.String(),
			MaxOpenPositions:    msg.Params.MaxOpenPositions,
			PoolOpenThreshold:   msg.Params.PoolOpenThreshold.String(),
			SafetyFactor:        msg.Params.SafetyFactor.String(),
			WhitelistingEnabled: msg.Params.WhitelistingEnabled,
			EpochLength:         msg.Params.EpochLength,
			FallbackEnabled:     msg.Params.FallbackEnabled,
			NumberPerBlock:      msg.Params.NumberPerBlock,
		},
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateParamsResponse{}, nil
}
