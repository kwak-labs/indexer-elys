package masterchef

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgAddExternalIncentive struct {
	Sender         string `json:"sender"`
	RewardDenom    string `json:"reward_denom"`
	PoolID         uint64 `json:"pool_id"`
	FromBlock      int64  `json:"from_block"`
	ToBlock        int64  `json:"to_block"`
	AmountPerBlock string `json:"amount_per_block"`
	TotalAmount    string `json:"total_amount"`
}

func (m MsgAddExternalIncentive) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
