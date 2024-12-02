package oracle

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type Params struct {
	BandChannelSource string        `json:"band_channel_source"`
	OracleScriptID    uint64        `json:"oracle_script_id"`
	Multiplier        uint64        `json:"multiplier"`
	AskCount          uint64        `json:"ask_count"`
	MinCount          uint64        `json:"min_count"`
	FeeLimit          []types.Token `json:"fee_limit"`
	PrepareGas        uint64        `json:"prepare_gas"`
	ExecuteGas        uint64        `json:"execute_gas"`
	ClientID          string        `json:"client_id"`
	BandEpoch         string        `json:"band_epoch"`
	PriceExpiryTime   uint64        `json:"price_expiry_time"`
	LifeTimeInBlocks  uint64        `json:"life_time_in_blocks"`
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
