package leveragelp

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type Params struct {
	LeverageMax         string `json:"leverage_max"`
	MaxOpenPositions    int64  `json:"max_open_positions"`
	PoolOpenThreshold   string `json:"pool_open_threshold"`
	SafetyFactor        string `json:"safety_factor"`
	WhitelistingEnabled bool   `json:"whitelisting_enabled"`
	EpochLength         int64  `json:"epoch_length"`
	FallbackEnabled     bool   `json:"fallback_enabled"`
	NumberPerBlock      int64  `json:"number_per_block"`
}

type MsgUpdateParams struct {
	Authority string `json:"authority"`
	Params    Params `json:"params"`
}

func (m MsgUpdateParams) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
