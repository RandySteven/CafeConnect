package midtrans_client

import (
	"context"
	"log"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func (m *midtransClient) CreateTransaction(ctx context.Context, request *MidtransRequest) (result *MidtransResponse, err error) {
	log.Println(`snap clientnya nil kah ?`, m.snapClient == nil)
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

	log.Println(response.Token)
	log.Println(response.RedirectURL)
	log.Println(response.StatusCode)

	result = &MidtransResponse{
		Token:       response.Token,
		RedirectURL: response.RedirectURL,
	}

	return result, nil
}
