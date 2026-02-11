package midtrans_usecases

import (
	"context"
	"fmt"

	"github.com/RandySteven/CafeConnect/be/entities/models"
)

func (m *midtransWorkflow) checkTransactionHeader(ctx context.Context, transactionCode string) (result *models.TransactionHeader, err error) {
	transactionHeader, err := m.transactionHeaderRepository.FindByTransactionCode(ctx, transactionCode)
	if err != nil {
		return nil, err
	}

	if transactionHeader.ID == 0 {
		return nil, fmt.Errorf("transaction header not found")
	}

	return transactionHeader, nil
}
