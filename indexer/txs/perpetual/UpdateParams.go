package perpetual

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type Params struct {
	LeverageMax                                    string `json:"leverage_max"`
	BorrowInterestRateMax                          string `json:"borrow_interest_rate_max"`
	BorrowInterestRateMin                          string `json:"borrow_interest_rate_min"`
	BorrowInterestRateIncrease                     string `json:"borrow_interest_rate_increase"`
	BorrowInterestRateDecrease                     string `json:"borrow_interest_rate_decrease"`
	HealthGainFactor                               string `json:"health_gain_factor"`
	MaxOpenPositions                               int64  `json:"max_open_positions"`
	PoolOpenThreshold                              string `json:"pool_open_threshold"`
	ForceCloseFundPercentage                       string `json:"force_close_fund_percentage"`
	ForceCloseFundAddress                          string `json:"force_close_fund_address"`
	IncrementalBorrowInterestPaymentFundPercentage string `json:"incremental_borrow_interest_payment_fund_percentage"`
	IncrementalBorrowInterestPaymentFundAddress    string `json:"incremental_borrow_interest_payment_fund_address"`
	SafetyFactor                                   string `json:"safety_factor"`
	IncrementalBorrowInterestPaymentEnabled        bool   `json:"incremental_borrow_interest_payment_enabled"`
	WhitelistingEnabled                            bool   `json:"whitelisting_enabled"`
	PerpetualSwapFee                               string `json:"perpetual_swap_fee"`
	MaxLimitOrder                                  int64  `json:"max_limit_order"`
	FixedFundingRate                               string `json:"fixed_funding_rate"`
	MinimumLongTakeProfitPriceRatio                string `json:"minimum_long_take_profit_price_ratio"`
	MaximumLongTakeProfitPriceRatio                string `json:"maximum_long_take_profit_price_ratio"`
	MaximumShortTakeProfitPriceRatio               string `json:"maximum_short_take_profit_price_ratio"`
	EnableTakeProfitCustodyLiabilities             bool   `json:"enable_take_profit_custody_liabilities"`
	WeightBreakingFeeFactor                        string `json:"weight_breaking_fee_factor"`
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
