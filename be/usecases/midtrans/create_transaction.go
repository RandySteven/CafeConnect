package midtrans_usecases

import (
	"context"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
)

func (m *midtransWorkflow) createMidtransTransaction(ctx context.Context, request *midtrans_client.MidtransRequest) (*midtrans_client.MidtransResponse, error) {
	midtransResponse, err := m.midtrans.CreateTransaction(ctx, request)
	if err != nil {
		return nil, err
	}

	_, err = m.midtransTransactionRepository.Save(ctx, &models.MidtransTransaction{
		TransactionCode: request.TransactionCode,
		TotalAmt:        request.GrossAmt,
		Token:           midtransResponse.Token,
		RedirectURL:     midtransResponse.RedirectURL})
	if err != nil {
		return nil, err
	}

	return midtransResponse, nil
}
