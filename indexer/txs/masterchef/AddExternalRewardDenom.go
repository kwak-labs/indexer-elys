package masterchef

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgAddExternalRewardDenom struct {
	Authority   string `json:"authority"`
	RewardDenom string `json:"reward_denom"`
	MinAmount   string `json:"min_amount"`
	Supported   bool   `json:"supported"`
}

func (m MsgAddExternalRewardDenom) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
