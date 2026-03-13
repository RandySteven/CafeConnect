package midtrans_usecases

import (
	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/midtrans/midtrans-go"
)

type (
	ExecutionData struct {
		Request *requests.CreateTransactionRequest

		Message           *messages.TransactionMidtransMessage
		MidtransRequest   *midtrans_client.MidtransRequest
		TransactionHeader *models.TransactionHeader
		Items             []midtrans.ItemDetails
		TotalAmount       int64
		MidtransResponse  *midtrans_client.MidtransResponse

		NextActivity    string `json:"next_activity,omitempty"`
		CurrentActivity string `json:"current_activity,omitempty"`
	}
)

func (m *ExecutionData) GetActivity() string     { return m.CurrentActivity }
func (m *ExecutionData) SetActivity(name string) { m.CurrentActivity = name }
