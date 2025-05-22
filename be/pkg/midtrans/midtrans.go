package midtrans_client

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type (
	MidtransRequest struct {
		FName string
		LName string
		Phone string

		Items []midtrans.ItemDetails
	}

	MidtransResponse struct {
		Token       string
		RedirectURL string
	}

	Midtrans interface {
		CreateTransaction(ctx context.Context, request *MidtransRequest)
	}

	midtransClient struct {
		snapClient *snap.Client
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

	sc := &snap.Client{
		ServerKey: midtransConf.ServerKey,
		Env:       midtransEnv,
	}

	return &midtransClient{
		snapClient: sc,
	}, nil
}
