package stablestake

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type Params struct {
	DepositDenom         string `json:"deposit_denom"`
	RedemptionRate       string `json:"redemption_rate"`
	EpochLength          int64  `json:"epoch_length"`
	InterestRate         string `json:"interest_rate"`
	InterestRateMax      string `json:"interest_rate_max"`
	InterestRateMin      string `json:"interest_rate_min"`
	InterestRateIncrease string `json:"interest_rate_increase"`
	InterestRateDecrease string `json:"interest_rate_decrease"`
	HealthGainFactor     string `json:"health_gain_factor"`
	TotalValue           string `json:"total_value"`
	MaxLeverageRatio     string `json:"max_leverage_ratio"`
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
