package auto_transfer_usecases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var (
	nonRetryableErrorTypes = []string{
		context.DeadlineExceeded.Error(),
		context.Canceled.Error(),
		sql.ErrConnDone.Error(),
	}
)

func (t *autoTransferWorkflow) autoTransferWorkflow(workflowCtx workflow.Context, userID uint64, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, err error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        5,
			InitialInterval:        1 * time.Minute,
			BackoffCoefficient:     2,
			MaximumInterval:        5 * time.Minute,
			NonRetryableErrorTypes: nonRetryableErrorTypes,
		},
	}
	ctx := workflow.WithActivityOptions(workflowCtx, ao)

	executionData := &ExecutionData{
		UserID:  userID,
		Request: request,
	}

	if err := t.workflow.Execute(ctx, executionData); err != nil {
		return nil, err
	}

	// If stock deduction failed and compensation was applied, stop here.
	if executionData.StockDeductionFailed {
		return nil, fmt.Errorf("stock deduction failed, compensation applied — deducted stock has been restored")
	}

	// Start MidtransTransaction child workflow and signal it,
	// matching the pattern in transactions/transaction.go.
	midtransResponse, err := t.midtransSignalTransaction(ctx, request, executionData)
	if err != nil {
		return nil, err
	}

	return &responses.TransactionReceiptResponse{
		ID:               executionData.TransactionHeader.ID,
		TransactionCode:  executionData.TransactionHeader.TransactionCode,
		Status:           executionData.TransactionHeader.Status,
		TransactionAt:    executionData.TransactionHeader.TransactionAt.Local(),
		MidtransResponse: midtransResponse,
	}, nil
}

func (t *autoTransferWorkflow) midtransSignalTransaction(ctx workflow.Context, request *requests.CreateTransactionRequest, executionData *ExecutionData) (*midtrans_client.MidtransResponse, error) {
	var midtransResponse *midtrans_client.MidtransResponse

	err := t.workflow.StartChildWorkflow(ctx, fmt.Sprintf("%s-%s", sgMidtrans, request.IdempotencyKey), "MidtransTransaction", executionData.MidtransMessage, &midtransResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to start midtrans child workflow: %w", err)
	}

	return midtransResponse, nil
}
