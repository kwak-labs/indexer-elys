package amm

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type SwapAmountInRoute struct {
	PoolID        uint64 `json:"pool_id"`
	TokenOutDenom string `json:"token_out_denom"`
}

type MsgSwapExactAmountIn struct {
	Sender            string              `json:"sender"`
	Routes            []SwapAmountInRoute `json:"routes"`
	TokenIn           types.Token         `json:"token_in"`
	TokenOutMinAmount string              `json:"token_out_min_amount"`
	Recipient         string              `json:"recipient"`
	SwapFee           string              `json:"swap_fee"`
	Discount          string              `json:"discount"`
	AmountOut         types.Token         `json:"amount_out"`
}

func (m MsgSwapExactAmountIn) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
