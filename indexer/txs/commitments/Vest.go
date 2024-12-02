package commitments

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgVest struct {
	Creator string      `json:"creator"`
	Token   types.Token `json:"token"`
}

func (m MsgVest) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
