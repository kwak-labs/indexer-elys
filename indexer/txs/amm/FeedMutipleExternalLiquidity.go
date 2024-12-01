package amm

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type AssetAmountDepth struct {
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
	Depth  string `json:"depth"`
}

type ExternalLiquidity struct {
	PoolID          uint64             `json:"pool_id"`
	AmountDepthInfo []AssetAmountDepth `json:"amount_depth_info"`
}

type MsgFeedMultipleExternalLiquidity struct {
	Sender    string              `json:"sender"`
	Liquidity []ExternalLiquidity `json:"liquidity"`
}

func (m MsgFeedMultipleExternalLiquidity) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
