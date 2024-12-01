package amm

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type PoolAsset struct {
	Token types.Token `json:"token"`
}

type PoolParams struct {
	SwapFee   string `json:"swap_fee"`
	UseOracle bool   `json:"use_oracle"`
	FeeDenom  string `json:"fee_denom"`
}

type MsgCreatePool struct {
	Sender     string      `json:"sender"`
	PoolParams PoolParams  `json:"pool_params"`
	PoolAssets []PoolAsset `json:"pool_assets"`
	PoolID     uint64      `json:"pool_id"`
}

func (m MsgCreatePool) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
