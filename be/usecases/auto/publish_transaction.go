package auto_transfer_usecases

import (
	"context"

	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/utils"
)

func (t *autoTransferWorkflow) publishTransaction(ctx context.Context, executionData *TransferExecutionData) (*TransferExecutionData, error) {
	fname, lname := utils.FirstLastName(executionData.User.Name)

	// Build the Midtrans signal message for the child workflow
	executionData.MidtransMessage = &messages.TransactionMidtransMessage{
		UserID:            executionData.User.ID,
		FName:             fname,
		Email:             executionData.User.Email,
		Phone:             executionData.User.PhoneNumber,
		LName:             lname,
		TransactionCode:   executionData.TransactionHeader.TransactionCode,
		CafeID:            executionData.Cafe.ID,
		CafeFranchiseName: executionData.Franchise.Name,
		CheckoutList:      executionData.Request.Checkouts,
	}

	return executionData, nil
}
