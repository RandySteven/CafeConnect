package auto_transfer_usecases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
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

	lao := workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        5,
			InitialInterval:        1 * time.Minute,
			BackoffCoefficient:     2,
			MaximumInterval:        5 * time.Minute,
			NonRetryableErrorTypes: nonRetryableErrorTypes,
		},
	}
	workflowCtx = workflow.WithLocalActivityOptions(workflowCtx, lao)

	stockDeductionActivityOptions := workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: 2 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        3,
			InitialInterval:        30 * time.Second,
		},
	}
	workflowDeductCtx := workflow.WithLocalActivityOptions(workflowCtx, stockDeductionActivityOptions)

	// Run checkUser and checkCafe in parallel
	var user *models.User
	var cafe *models.Cafe

	userFuture := workflow.ExecuteLocalActivity(workflowCtx, t.checkUser, userID)
	cafeFuture := workflow.ExecuteLocalActivity(workflowCtx, t.checkCafe, request.CafeID)

	if err := userFuture.Get(workflowCtx, &user); err != nil {
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	if err := cafeFuture.Get(workflowCtx, &cafe); err != nil {
		return nil, fmt.Errorf("failed to check cafe: %w", err)
	}

	// Check franchise
	var franchise *models.CafeFranchise
	if err := workflow.ExecuteLocalActivity(workflowCtx, t.checkFranchise, cafe.CafeFranchiseID).Get(workflowCtx, &franchise); err != nil {
		return nil, err
	}

	// Save transaction header
	var transactionHeader *models.TransactionHeader
	if err := workflow.ExecuteLocalActivity(workflowCtx, t.saveTransactionHeader, user.ID, request).Get(workflowCtx, &transactionHeader); err != nil {
		return nil, err
	}

	// Publish transaction
	if err := workflow.ExecuteLocalActivity(workflowCtx, t.publishTransaction, user, cafe, franchise, transactionHeader, request).Get(workflowCtx, nil); err != nil {
		return nil, err
	}

	// Stock deduction
	if err := workflow.ExecuteLocalActivity(workflowDeductCtx, t.stockDeduction, request.Checkouts).Get(workflowCtx, nil); err != nil {
		return nil, err
	}

	// Save transaction detail
	if err := workflow.ExecuteLocalActivity(workflowCtx, t.saveTransactionDetail, transactionHeader.ID, request.Checkouts).Get(workflowCtx, nil); err != nil {
		return nil, err
	}

	return &responses.TransactionReceiptResponse{
		ID:              transactionHeader.ID,
		TransactionCode: transactionHeader.TransactionCode,
		Status:          transactionHeader.Status,
		TransactionAt:   transactionHeader.TransactionAt.Local(),
	}, nil
}
