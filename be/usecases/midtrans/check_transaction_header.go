package midtrans_usecases

import (
	"context"
	"fmt"
)

func (m *midtransWorkflow) checkTransactionHeader(ctx context.Context, executionData *MidtransExecutionData) (*MidtransExecutionData, error) {
	transactionHeader, err := m.transactionHeaderRepository.FindByTransactionCode(ctx, executionData.Message.TransactionCode)
	if err != nil {
		return nil, err
	}

	if transactionHeader.ID == 0 {
		return nil, fmt.Errorf("transaction header not found")
	}

	executionData.TransactionHeader = transactionHeader
	return executionData, nil
}
