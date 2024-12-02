package tokenomics

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

// type InflationEntry struct {
// 	LmRewards         string `json:"lm_rewards"`
// 	IcsStakingRewards string `json:"ics_staking_rewards"`
// 	CommunityFund     string `json:"community_fund"`
// 	StrategicReserve  string `json:"strategic_reserve"`
// 	TeamTokensVested  string `json:"team_tokens_vested"`
// }

type MsgCreateTimeBasedInflation struct {
	Authority        string         `json:"authority"`
	StartBlockHeight uint64         `json:"start_block_height"`
	EndBlockHeight   uint64         `json:"end_block_height"`
	Description      string         `json:"description"`
	Inflation        InflationEntry `json:"inflation"`
}

func (m MsgCreateTimeBasedInflation) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
