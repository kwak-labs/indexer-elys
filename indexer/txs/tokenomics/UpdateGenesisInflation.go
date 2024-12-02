package tokenomics

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type InflationEntry struct {
	LmRewards         uint64 `json:"lm_rewards"`
	IcsStakingRewards uint64 `json:"ics_staking_rewards"`
	CommunityFund     uint64 `json:"community_fund"`
	StrategicReserve  uint64 `json:"strategic_reserve"`
	TeamTokensVested  uint64 `json:"team_tokens_vested"`
}

type MsgUpdateGenesisInflation struct {
	Authority             string          `json:"authority"`
	Inflation             *InflationEntry `json:"inflation"`
	SeedVesting           uint64          `json:"seed_vesting"`
	StrategicSalesVesting uint64          `json:"strategic_sales_vesting"`
}

func (m MsgUpdateGenesisInflation) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
	mergedData := types.GenericTransaction{
		BaseTransaction: transaction,
		Data:            m,
	}

	err := database.ProcessNewTx(mergedData, transaction.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing transaction: %w", err)
	}

	return types.Response{}, nil
}
