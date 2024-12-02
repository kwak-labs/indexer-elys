package amm

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

// Defined in SwapExactAmountIn.go
// type SwapAmountInRoute struct {
// 	PoolId        uint64 `json:"pool_id"`
// 	TokenOutDenom string `json:"token_out_denom"`
// }

// Defined in SwapExactAmountOut.go
// type SwapAmountOutRoute struct {
// 	PoolId       uint64 `json:"pool_id"`
// 	TokenInDenom string `json:"token_in_denom"`
// }

type MsgSwapByDenom struct {
	Sender    string               `json:"sender"`
	Amount    types.Token          `json:"amount"`
	MinAmount types.Token          `json:"min_amount"`
	MaxAmount types.Token          `json:"max_amount"`
	DenomIn   string               `json:"denom_in"`
	DenomOut  string               `json:"denom_out"`
	Recipient string               `json:"recipient"`
	InRoute   []SwapAmountInRoute  `json:"in_route,omitempty"`
	OutRoute  []SwapAmountOutRoute `json:"out_route,omitempty"`
	SpotPrice string               `json:"spot_price"`
	SwapFee   string               `json:"swap_fee"`
	Discount  string               `json:"discount"`
	TokenOut  types.Token          `json:"token_out"`
}

func (m MsgSwapByDenom) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
