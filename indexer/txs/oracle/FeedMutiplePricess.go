package oracle

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type FeedPrice struct {
	Asset  string `json:"asset"`
	Price  string `json:"price"`
	Source string `json:"source"`
}

type MsgFeedMultiplePrices struct {
	Creator    string      `json:"creator"`
	FeedPrices []FeedPrice `json:"feed_prices"`
	Timestamp  uint64      `json:"timestamp"`
}

func (m MsgFeedMultiplePrices) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
