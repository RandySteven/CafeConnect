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

	return nil
}
