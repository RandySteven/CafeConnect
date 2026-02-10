package transactions_usecases

import (
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"go.temporal.io/sdk/workflow"
)

func (t *transactionWorkflow) transactionWorkflow(workflowCtx workflow.Context, userID uint64, request *requests.CreateTransactionRequest) (*responses.TransactionReceiptResponse, error) {
	lao := workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Second,
	}
	workflowCtx = workflow.WithLocalActivityOptions(workflowCtx, lao)

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
		return nil, fmt.Errorf("failed to check franchise: %w", err)
	}

	// Save transaction header
	var transactionHeader *models.TransactionHeader
	if err := workflow.ExecuteLocalActivity(workflowCtx, t.saveTransactionHeader, user.ID, request).Get(workflowCtx, &transactionHeader); err != nil {
		return nil, fmt.Errorf("failed to save transaction header: %w", err)
	}

	// Publish transaction
	if err := workflow.ExecuteLocalActivity(workflowCtx, t.publishTransaction, user, cafe, franchise, transactionHeader, request).Get(workflowCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to publish transaction: %w", err)
	}

	return &responses.TransactionReceiptResponse{
		ID:              transactionHeader.ID,
		TransactionCode: transactionHeader.TransactionCode,
		Status:          transactionHeader.Status,
		TransactionAt:   transactionHeader.TransactionAt.Local(),
	}, nil
}
