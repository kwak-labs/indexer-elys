package estaking

import (
	"github.com/elys-network/elys/indexer/types"
)

// UpdateParams Message
type MsgUpdateParams struct {
	Authority string      `json:"authority"`
	Params    interface{} `json:"params"`
}

func (m MsgUpdateParams) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
	mergedData := types.GenericTransaction{
		BaseTransaction: transaction,
		Data:            m,
	}
	return types.Response{}, database.ProcessNewTx(mergedData, transaction.Author)
}

// WithdrawReward Message
type MsgWithdrawReward struct {
	DelegatorAddress string        `json:"delegator_address"`
	ValidatorAddress string        `json:"validator_address"`
	Amount           []types.Token `json:"amount"`
}

func (m MsgWithdrawReward) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
	mergedData := types.GenericTransaction{
		BaseTransaction: transaction,
		Data:            m,
	}
	return types.Response{}, database.ProcessNewTx(mergedData, transaction.Author)
}

// WithdrawElysStakingRewards Message
type MsgWithdrawElysStakingRewards struct {
	DelegatorAddress string        `json:"delegator_address"`
	Validators       []string      `json:"validators"`
	Amount           []types.Token `json:"amount"`
}

func (m MsgWithdrawElysStakingRewards) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
	mergedData := types.GenericTransaction{
		BaseTransaction: transaction,
		Data:            m,
	}
	return types.Response{}, database.ProcessNewTx(mergedData, transaction.Author)
}

// WithdrawAllRewards Message
type MsgWithdrawAllRewards struct {
	DelegatorAddress string        `json:"delegator_address"`
	Validators       []string      `json:"validators"`
	Amount           []types.Token `json:"amount"`
}

func (m MsgWithdrawAllRewards) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
	mergedData := types.GenericTransaction{
		BaseTransaction: transaction,
		Data:            m,
	}
	return types.Response{}, database.ProcessNewTx(mergedData, transaction.Author)
}
