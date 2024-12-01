package oracle

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

// Defined in FeedMutiplePrices.go
// type FeedPrice struct {
// 	Asset  string `json:"asset"`
// 	Price  string `json:"price"`
// 	Source string `json:"source"`
// }

type MsgFeedPrice struct {
	Provider    string `json:"provider"`
	Asset       string `json:"asset"`
	Price       string `json:"price"`
	Source      string `json:"source"`
	Timestamp   uint64 `json:"timestamp"`
	BlockHeight uint64 `json:"block_height"`
}

func (m MsgFeedPrice) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
