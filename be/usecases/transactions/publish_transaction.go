package transactions_usecases

import (
	"context"
	"fmt"

	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/utils"
)

func (t *transactionWorkflow) publishTransaction(ctx context.Context, user *models.User, cafe *models.Cafe, franchise *models.CafeFranchise, transactionHeader *models.TransactionHeader, request *requests.CreateTransactionRequest) error {
	fname, lname := utils.FirstLastName(user.Name)
	err := t.transactionTopic.WriteMessage(ctx, utils.WriteJSONObject[messages.TransactionMidtransMessage](&messages.TransactionMidtransMessage{
		UserID:            user.ID,
		FName:             fname,
		Email:             user.Email,
		Phone:             user.PhoneNumber,
		LName:             lname,
		TransactionCode:   transactionHeader.TransactionCode,
		CafeID:            cafe.ID,
		CafeFranchiseName: franchise.Name,
		CheckoutList:      request.Checkouts,
	}))
	if err != nil {
		return fmt.Errorf("failed to publish transaction: %w", err)
	}
	// workflowInfo, err := t.workflow.GetWorkflowInfo(workflowCtx)
	// if err != nil {
	// 	return fmt.Errorf("failed to get workflow info: %w", err)
	// }
	// runID := workflowInfo.RunID
	// workflowID := "MidtransTransaction"
	// runID := t.workflow.GetWorkflowRunID()
	// err = t.workflow.SignalWorkflow(ctx, workflowID, runID, "MidtransTransaction", &messages.TransactionMidtransMessage{
	// 	TransactionCode: transactionHeader.TransactionCode,
	// 	CheckoutList:    request.Checkouts,
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to signal workflow: %w", err)
	// }

	// var midtransResponse *midtrans_client.MidtransResponse
	// queryResult, err := t.workflow.QueryWorkflow(ctx, workflowRun.GetID(), "MidtransTransaction", workflowRun.GetRunID())
	// if err != nil {
	// 	return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to query workflow`, err)
	// }

	// if queryResult == nil {
	// 	return nil, apperror.NewCustomError(apperror.ErrNotFound, `query result not found`, nil)
	// }

	// midtransResponse = queryResult.(*midtrans_client.MidtransResponse)
	// result.MidtransResponse = midtransResponse
	// return errors.New("mock error")
	return nil
}
