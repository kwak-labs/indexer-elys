package assetprofile

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgAddEntry struct {
	Authority                string   `json:"authority"`
	BaseDenom                string   `json:"base_denom"`
	Decimals                 uint64   `json:"decimals"`
	Denom                    string   `json:"denom"`
	Path                     string   `json:"path"`
	IbcChannelId             string   `json:"ibc_channel_id"`
	IbcCounterpartyChannelId string   `json:"ibc_counterparty_channel_id"`
	DisplayName              string   `json:"display_name"`
	DisplaySymbol            string   `json:"display_symbol"`
	Network                  string   `json:"network"`
	Address                  string   `json:"address"`
	ExternalSymbol           string   `json:"external_symbol"`
	TransferLimit            string   `json:"transfer_limit"`
	Permissions              []string `json:"permissions"`
	UnitDenom                string   `json:"unit_denom"`
	IbcCounterpartyDenom     string   `json:"ibc_counterparty_denom"`
	IbcCounterpartyChainId   string   `json:"ibc_counterparty_chain_id"`
	CommitEnabled            bool     `json:"commit_enabled"`
	WithdrawEnabled          bool     `json:"withdraw_enabled"`
}

func (m MsgAddEntry) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
