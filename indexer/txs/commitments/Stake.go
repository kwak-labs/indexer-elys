package commitments

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgStake struct {
	Creator           string      `json:"creator"`
	Token             types.Token `json:"token"`
	ValidatorAddresss string      `json:"validator_address"`
}

func (m MsgStake) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
	mergedData := types.GenericTransaction{
		BaseTransaction: transaction,
		Data:            m,
	}

	err := database.ProcessNewTx(mergedData, transaction.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing transaction: %w", err)
	}

	fmt.Println("Successfully Stored")

	return types.Response{}, nil
}
