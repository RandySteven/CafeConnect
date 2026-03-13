package midtrans_usecases

import (
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/messages"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"go.temporal.io/sdk/workflow"
)

func (m *midtransWorkflow) midtransTransaction(workflowCtx workflow.Context) (*midtrans_client.MidtransResponse, error) {
	// Wait for signal from the handler with the transaction data
	var message messages.TransactionMidtransMessage
	err := m.workflow.GetSignalResult(workflowCtx, "MidtransTransaction", &message)
	if err != nil {
		return nil, err
	}

	lao := workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Second,
	}
	workflowCtx = workflow.WithLocalActivityOptions(workflowCtx, lao)

	executionData := &ExecutionData{
		Message: &message,
	}

	err = m.workflow.Execute(workflowCtx, executionData)
	if err != nil {
		return nil, err
	}

	return executionData.MidtransResponse, nil
}
