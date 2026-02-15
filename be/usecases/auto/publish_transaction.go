package auto_transfer_usecases

import (
	"context"

	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/utils"
)

func (t *autoTransferWorkflow) publishTransaction(ctx context.Context, state *TransferState) (*TransferState, error) {
	fname, lname := utils.FirstLastName(state.User.Name)

	// Build the Midtrans signal message for the child workflow
	state.MidtransMessage = &messages.TransactionMidtransMessage{
		UserID:            state.User.ID,
		FName:             fname,
		Email:             state.User.Email,
		Phone:             state.User.PhoneNumber,
		LName:             lname,
		TransactionCode:   state.TransactionHeader.TransactionCode,
		CafeID:            state.Cafe.ID,
		CafeFranchiseName: state.Franchise.Name,
		CheckoutList:      state.Request.Checkouts,
	}

	return state, nil
}
