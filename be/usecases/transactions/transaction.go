package transactions_usecases

import (
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/utils"
	"go.temporal.io/sdk/workflow"
)

// transactionResult is the internal result of the transaction workflow.
// It carries both the receipt and the data needed to signal the Midtrans workflow.
type transactionResult struct {
	Receipt           *responses.TransactionReceiptResponse `json:"receipt"`
	UserID            uint64                                `json:"user_id"`
	FName             string                                `json:"f_name"`
	LName             string                                `json:"l_name"`
	Email             string                                `json:"email"`
	Phone             string                                `json:"phone"`
	CafeID            uint64                                `json:"cafe_id"`
	CafeFranchiseName string                                `json:"cafe_franchise_name"`
}

func (t *transactionWorkflow) transactionWorkflow(workflowCtx workflow.Context, userID uint64, request *requests.CreateTransactionRequest) (*transactionResult, error) {
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

	fname, lname := utils.FirstLastName(user.Name)

	return &transactionResult{
		Receipt: &responses.TransactionReceiptResponse{
			ID:              transactionHeader.ID,
			TransactionCode: transactionHeader.TransactionCode,
			Status:          transactionHeader.Status,
			TransactionAt:   transactionHeader.TransactionAt.Local(),
		},
		UserID:            user.ID,
		FName:             fname,
		LName:             lname,
		Email:             user.Email,
		Phone:             user.PhoneNumber,
		CafeID:            cafe.ID,
		CafeFranchiseName: franchise.Name,
	}, nil
}
