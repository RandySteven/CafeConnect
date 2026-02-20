package midtrans_usecases

import (
	"context"

	"github.com/RandySteven/CafeConnect/be/entities/models"
)

func (m *midtransWorkflow) createMidtransTransaction(ctx context.Context, executionData *MidtransExecutionData) (*MidtransExecutionData, error) {
	midtransResponse, err := m.midtrans.CreateTransaction(ctx, executionData.MidtransRequest)
	if err != nil {
		return nil, err
	}

	_, err = m.midtransTransactionRepository.Save(ctx, &models.MidtransTransaction{
		TransactionCode: executionData.Message.TransactionCode,
		TotalAmt:        executionData.TotalAmount,
		Token:           midtransResponse.Token,
		RedirectURL:     midtransResponse.RedirectURL})
	if err != nil {
		return nil, err
	}

	executionData.MidtransResponse = midtransResponse
	return executionData, nil
}
