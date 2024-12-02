package masterchef

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type IncentiveInfo struct {
	EdenAmountPerYear string `json:"eden_amount_per_year"`
	BlocksDistributed int64  `json:"blocks_distributed"`
}

type SupportedRewardDenom struct {
	Denom     string `json:"denom"`
	MinAmount string `json:"min_amount"`
}

type Params struct {
	LpIncentives            *IncentiveInfo         `json:"lp_incentives,omitempty"`
	RewardPortionForLps     string                 `json:"reward_portion_for_lps"`
	RewardPortionForStakers string                 `json:"reward_portion_for_stakers"`
	MaxEdenRewardAprLps     string                 `json:"max_eden_reward_apr_lps"`
	SupportedRewardDenoms   []SupportedRewardDenom `json:"supported_reward_denoms"`
	ProtocolRevenueAddress  string                 `json:"protocol_revenue_address"`
}

type MsgUpdateParams struct {
	Authority string `json:"authority"`
	Params    Params `json:"params"`
}

func (m MsgUpdateParams) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
