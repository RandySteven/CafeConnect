package midtrans_usecases

import (
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"go.temporal.io/sdk/workflow"
)

func (m *midtransWorkflow) midtransTransaction(workflowCtx workflow.Context) (*midtrans_client.MidtransResponse, error) {
	// Wait for signal from the handler with the transaction data
	var message messages.TransactionMidtransMessage
	signalChan := workflow.GetSignalChannel(workflowCtx, "MidtransTransaction")
	signalChan.Receive(workflowCtx, &message)

	lao := workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Second,
	}
	workflowCtx = workflow.WithLocalActivityOptions(workflowCtx, lao)

	var transactionHeader *models.TransactionHeader
	if err := workflow.ExecuteLocalActivity(workflowCtx, m.checkTransactionHeader, message.TransactionCode).Get(workflowCtx, &transactionHeader); err != nil {
		return nil, fmt.Errorf("failed to check transaction header: %w", err)
	}

	var checkoutList *midtransCheckOut
	if err := workflow.ExecuteLocalActivity(workflowCtx, m.checkoutList, message.CafeFranchiseName, message.CheckoutList).Get(workflowCtx, &checkoutList); err != nil {
		return nil, fmt.Errorf("failed to checkout list: %w", err)
	}

	var midtransResponse *midtrans_client.MidtransResponse
	if err := workflow.ExecuteLocalActivity(workflowCtx, m.createMidtransTransaction, &midtrans_client.MidtransRequest{
		FName:           message.FName,
		LName:           message.LName,
		Email:           message.Email,
		Phone:           message.Phone,
		TransactionCode: message.TransactionCode,
		GrossAmt:        checkoutList.TotalAmount,
		Items:           checkoutList.Items,
	}).Get(workflowCtx, &midtransResponse); err != nil {
		return nil, fmt.Errorf("failed to create midtrans transaction: %w", err)
	}

	return midtransResponse, nil
}
