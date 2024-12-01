package leveragelp

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type PositionOpen struct {
	ID          uint64      `json:"id"`
	Address     string      `json:"address"`
	Collateral  types.Token `json:"collateral"`
	Liabilities string      `json:"liabilities"`
	Health      string      `json:"health"`
}

type MsgOpen struct {
	Creator          string       `json:"creator"`
	CollateralAsset  string       `json:"collateral_asset"`
	CollateralAmount string       `json:"collateral_amount"`
	AmmPoolID        uint64       `json:"amm_pool_id"`
	Leverage         string       `json:"leverage"`
	StopLossPrice    string       `json:"stop_loss_price"`
	Position         PositionOpen `json:"position"`
}

func (m MsgOpen) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
