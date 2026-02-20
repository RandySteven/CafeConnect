package midtrans_usecases

import (
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/midtrans/midtrans-go"
	"go.temporal.io/sdk/workflow"
)

type (
	MidtransExecutionData struct {
		Request *requests.CreateTransactionRequest

		Message           *messages.TransactionMidtransMessage
		MidtransRequest   *midtrans_client.MidtransRequest
		TransactionHeader *models.TransactionHeader
		Items             []midtrans.ItemDetails
		TotalAmount       int64
		MidtransResponse  *midtrans_client.MidtransResponse

		NextActivity string `json:"next_activity,omitempty"`
	}
)

func (m *MidtransExecutionData) GetNextActivity() string     { return m.NextActivity }
func (m *MidtransExecutionData) SetNextActivity(name string) { m.NextActivity = name }

func (m *midtransWorkflow) midtransTransaction(workflowCtx workflow.Context) (*midtrans_client.MidtransResponse, error) {
	// Wait for signal from the handler with the transaction data
	var message messages.TransactionMidtransMessage
	signalChan := workflow.GetSignalChannel(workflowCtx, "MidtransTransaction")
	signalChan.Receive(workflowCtx, &message)

	lao := workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Second,
	}
	workflowCtx = workflow.WithLocalActivityOptions(workflowCtx, lao)

	executionData := &MidtransExecutionData{
		Message: &message,
	}

	err := m.workflow.Execute(workflowCtx, executionData)
	if err != nil {
		return nil, err
	}

	return executionData.MidtransResponse, nil
}
