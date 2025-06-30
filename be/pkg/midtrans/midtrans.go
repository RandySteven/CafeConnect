package midtrans_client

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type (
	MidtransRequest struct {
		FName           string
		LName           string
		Email           string
		Phone           string
		TransactionCode string
		GrossAmt        int64
		Items           []midtrans.ItemDetails
	}

	MidtransResponse struct {
		Token       string `json:"token"`
		RedirectURL string `json:"redirect_url"`
	}

	Midtrans interface {
		CreateTransaction(ctx context.Context, request *MidtransRequest) (result *MidtransResponse, err error)
		CheckTransaction(ctx context.Context, orderId string) (response *coreapi.TransactionStatusResponse, err error)
		CheckTransactionHistory(ctx context.Context, orderId string)
	}

	midtransClient struct {
		snapClient *snap.Client
		coreApi    *coreapi.Client
	}
)

func NewMidtrans(config *configs.Config) (*midtransClient, error) {
	midtransConf := config.Config.Midtrans
	var midtransEnv midtrans.EnvironmentType

	if midtransConf.Environment == "SANDBOX" || midtransConf.Environment == "" {
		midtransEnv = midtrans.Sandbox
	} else {
		midtransEnv = midtrans.Production
	}

	sc := &snap.Client{}
	sc.New(midtransConf.ServerKey, midtransEnv)

	ca := &coreapi.Client{}
	ca.New(midtransConf.ServerKey, midtransEnv)

	return &midtransClient{
		snapClient: sc,
		coreApi:    ca,
	}, nil
}
