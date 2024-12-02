package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerTokenomicsTypes "github.com/elys-network/elys/indexer/txs/tokenomics"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func (k msgServer) UpdateGenesisInflation(goCtx context.Context, msg *types.MsgUpdateGenesisInflation) (*types.MsgUpdateGenesisInflationResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	genesisInflation := types.GenesisInflation{
		Authority:             msg.Authority,
		Inflation:             msg.Inflation,
		SeedVesting:           msg.SeedVesting,
		StrategicSalesVesting: msg.StrategicSalesVesting,
	}

	k.SetGenesisInflation(ctx, genesisInflation)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerTokenomicsTypes.MsgUpdateGenesisInflation{
		Authority: msg.Authority,
		Inflation: &indexerTokenomicsTypes.InflationEntry{
			LmRewards:         msg.Inflation.LmRewards,
			IcsStakingRewards: msg.Inflation.IcsStakingRewards,
			CommunityFund:     msg.Inflation.CommunityFund,
			StrategicReserve:  msg.Inflation.StrategicReserve,
			TeamTokensVested:  msg.Inflation.TeamTokensVested,
		},
		SeedVesting:           msg.SeedVesting,
		StrategicSalesVesting: msg.StrategicSalesVesting,
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateGenesisInflationResponse{}, nil
}
