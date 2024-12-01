package amm

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgExitPool struct {
	Sender        string        `json:"sender"`
	PoolID        uint64        `json:"pool_id"`
	MinAmountsOut []types.Token `json:"min_amounts_out"`
	ShareAmountIn string        `json:"share_amount_in"`
	TokenOutDenom string        `json:"token_out_denom"`
	TokenOut      []types.Token `json:"token_out"`
}

func (m MsgExitPool) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
