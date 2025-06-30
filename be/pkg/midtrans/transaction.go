package midtrans_client

import (
	"context"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

func (m *midtransClient) CreateTransaction(ctx context.Context, request *MidtransRequest) (result *MidtransResponse, err error) {
	response, midtransErr := m.snapClient.CreateTransaction(
		&snap.Request{
			CustomerDetail: &midtrans.CustomerDetails{
				FName: request.FName,
				LName: request.LName,
				Phone: request.Phone,
				Email: request.Email,
			},
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  request.TransactionCode,
				GrossAmt: request.GrossAmt,
			},
			Items: &request.Items,
		},
	)
	if midtransErr != nil {
		return nil, midtransErr
	}

	result = &MidtransResponse{
		Token:       response.Token,
		RedirectURL: response.RedirectURL,
	}

	return result, nil
}

func (m *midtransClient) CheckTransaction(ctx context.Context, orderId string) (response *coreapi.TransactionStatusResponse, err error) {
	response, midtransErr := m.coreApi.CheckTransaction(orderId)
	if midtransErr != nil {
		return nil, midtransErr
	}
	return response, nil
}

func (m *midtransClient) CheckTransactionHistory(ctx context.Context, orderId string) {
}
