package amm

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgJoinPool struct {
	Sender         string        `json:"sender"`
	PoolID         uint64        `json:"pool_id"`
	MaxAmountsIn   []types.Token `json:"max_amounts_in"`
	ShareAmountOut string        `json:"share_amount_out"`
	TokenIn        []types.Token `json:"token_in"` // Actual tokens used to join the pool
}

func (m MsgJoinPool) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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