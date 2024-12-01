package amm

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type Params struct {
	PoolCreationFee             string   `json:"pool_creation_fee"`
	SlippageTrackDuration       uint64   `json:"slippage_track_duration"`
	BaseAssets                  []string `json:"base_assets"`
	WeightBreakingFeeExponent   string   `json:"weight_breaking_fee_exponent"`
	WeightBreakingFeeMultiplier string   `json:"weight_breaking_fee_multiplier"`
	WeightBreakingFeePortion    string   `json:"weight_breaking_fee_portion"`
	WeightRecoveryFeePortion    string   `json:"weight_recovery_fee_portion"`
	ThresholdWeightDifference   string   `json:"threshold_weight_difference"`
	AllowedPoolCreators         []string `json:"allowed_pool_creators"`
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
