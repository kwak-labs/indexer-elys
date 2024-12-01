package perpetual

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgClose struct {
	Creator     string      `json:"creator"`
	Id          uint64      `json:"id"`
	Amount      string      `json:"amount"`
	RepayAmount string      `json:"repay_amount"`
	Position    string      `json:"position"`
	Collateral  types.Token `json:"collateral"`
	Custody     types.Token `json:"custody"`
	Liabilities types.Token `json:"liabilities"`
	// Profit/Loss tracking
	InitialValue   string `json:"initial_value"`
	FinalValue     string `json:"final_value"`
	ProfitLoss     string `json:"profit_loss"`
	ProfitLossPerc string `json:"profit_loss_perc"`
	// Additional MTP info
	CollateralAsset  string `json:"collateral_asset"`
	TradingAsset     string `json:"trading_asset"`
	LiabilitiesAsset string `json:"liabilities_asset"`
	MtpHealth        string `json:"mtp_health"`
	OpenPrice        string `json:"open_price"`
}

func (m MsgClose) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
