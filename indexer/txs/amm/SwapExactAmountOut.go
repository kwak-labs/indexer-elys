package amm

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type SwapAmountOutRoute struct {
	PoolID       uint64 `json:"pool_id"`
	TokenInDenom string `json:"token_in_denom"`
}

type MsgSwapExactAmountOut struct {
	Sender           string               `json:"sender"`
	Routes           []SwapAmountOutRoute `json:"routes"`
	TokenOut         types.Token          `json:"token_out"`
	TokenInMaxAmount string               `json:"token_in_max_amount"`
	Recipient        string               `json:"recipient"`
	TokenInAmount    types.Token          `json:"token_in_amount"`
	SwapFee          types.Token          `json:"swap_fee"`
	Discount         types.Token          `json:"discount"`
}

func (m MsgSwapExactAmountOut) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
