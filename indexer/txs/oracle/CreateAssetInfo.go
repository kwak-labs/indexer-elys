package oracle

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgCreateAssetInfo struct {
	Creator    string `json:"creator"`
	Denom      string `json:"denom"`
	Display    string `json:"display"`
	BandTicker string `json:"band_ticker"`
	ElysTicker string `json:"elys_ticker"`
	Decimal    uint64 `json:"decimal"`
}

func (m MsgCreateAssetInfo) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
