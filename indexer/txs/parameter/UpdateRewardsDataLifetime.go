package parameter

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgUpdateRewardsDataLifetime struct {
	Creator             string `json:"creator"`
	RewardsDataLifetime uint64 `json:"rewards_data_lifetime"`
}

func (m MsgUpdateRewardsDataLifetime) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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