package masterchef

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type PoolMultiplier struct {
	PoolID     uint64 `json:"pool_id"`
	Multiplier string `json:"multiplier"`
}

type MsgUpdatePoolMultipliers struct {
	Authority       string           `json:"authority"`
	PoolMultipliers []PoolMultiplier `json:"pool_multipliers"`
}

func (m MsgUpdatePoolMultipliers) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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