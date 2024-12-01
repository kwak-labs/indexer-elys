package amm

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

// Defined In CreatePool.go

// type PoolParams struct {
// 	SwapFee   string `json:"swap_fee"`
// 	UseOracle bool   `json:"use_oracle"`
// 	FeeDenom  string `json:"fee_denom"`
// }

type MsgUpdatePoolParams struct {
	Authority  string     `json:"authority"`
	PoolID     uint64     `json:"pool_id"`
	PoolParams PoolParams `json:"pool_params"`
}

func (m MsgUpdatePoolParams) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
