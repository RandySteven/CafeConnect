package auto_transfer_usecases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
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

// DeductedProduct tracks a stock deduction that was applied, so it can
// be rolled back by the restoreStock compensation activity.
type DeductedProduct struct {
	CafeProductID uint64 `json:"cafe_product_id"`
	Qty           uint64 `json:"qty"`
}

// TransferState is the shared serializable state threaded through all activities
// in the auto transfer pipeline. Each activity reads its inputs from state and
// writes its outputs back.
//
// It implements temporal_client.Navigable so activities can control branching
// by setting NextActivity (e.g., to trigger compensation on failure).
type TransferState struct {
	// Inputs
	UserID  uint64                             `json:"user_id"`
	Request *requests.CreateTransactionRequest `json:"request"`

	// Intermediate results — populated by activities
	User              *models.User                        `json:"user,omitempty"`
	Cafe              *models.Cafe                        `json:"cafe,omitempty"`
	Franchise         *models.CafeFranchise               `json:"franchise,omitempty"`
	TransactionHeader *models.TransactionHeader            `json:"transaction_header,omitempty"`
	MidtransMessage   *messages.TransactionMidtransMessage `json:"midtrans_message,omitempty"`
	MidtransResponse  *midtrans_client.MidtransResponse    `json:"midtrans_response,omitempty"`

	// Branching — controls which activity runs next.
	// If empty, Execute follows the default sequential order.
	// Set to an activity name to branch (e.g., for compensation).
	NextActivity string `json:"next_activity,omitempty"`

	// Compensation tracking
	DeductedProducts     []*DeductedProduct `json:"deducted_products,omitempty"`
	StockDeductionFailed bool               `json:"stock_deduction_failed,omitempty"`
}

// GetNextActivity implements temporal_client.Navigable.
func (s *TransferState) GetNextActivity() string { return s.NextActivity }

// SetNextActivity implements temporal_client.Navigable.
func (s *TransferState) SetNextActivity(name string) { s.NextActivity = name }

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

	state := &TransferState{
		UserID:  userID,
		Request: request,
	}

	if err := t.workflow.Execute(ctx, state); err != nil {
		return nil, err
	}

	// If stock deduction failed and compensation was applied, stop here.
	if state.StockDeductionFailed {
		return nil, fmt.Errorf("stock deduction failed, compensation applied — deducted stock has been restored")
	}

	// Start MidtransTransaction child workflow and signal it,
	// matching the pattern in transactions/transaction.go.
	childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: fmt.Sprintf("MidtransTransaction-%s", request.IdempotencyKey),
	})

	midtransRun := workflow.ExecuteChildWorkflow(childCtx, "MidtransTransaction")

	var childExecution workflow.Execution
	if err := midtransRun.GetChildWorkflowExecution().Get(ctx, &childExecution); err != nil {
		return nil, fmt.Errorf("failed to start midtrans child workflow: %w", err)
	}

	// Signal the child workflow with the prepared message from publishTransaction activity
	sigFuture := workflow.SignalExternalWorkflow(ctx, childExecution.ID, childExecution.RunID, "MidtransTransaction", state.MidtransMessage)
	if err := sigFuture.Get(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to signal midtrans workflow: %w", err)
	}

	// Wait for the child workflow to complete and get the Midtrans response
	var midtransResponse *midtrans_client.MidtransResponse
	if err := midtransRun.Get(childCtx, &midtransResponse); err != nil {
		return nil, fmt.Errorf("failed to get midtrans response: %w", err)
	}

	return &responses.TransactionReceiptResponse{
		ID:              state.TransactionHeader.ID,
		TransactionCode: state.TransactionHeader.TransactionCode,
		Status:          state.TransactionHeader.Status,
		TransactionAt:   state.TransactionHeader.TransactionAt.Local(),
		MidtransResponse: midtransResponse,
	}, nil
}
