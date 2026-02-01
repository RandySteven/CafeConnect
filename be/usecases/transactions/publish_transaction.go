package transactions_usecases

import (
	"context"

	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/utils"
)

func (t *transactionWorkflow) publishTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (err error) {
	transactionHeader := ctx.Value("transactionHeader").(*models.TransactionHeader)
	user := ctx.Value("user").(*models.User)
	cafe := ctx.Value("cafe").(*models.Cafe)
	cafeFranchise := ctx.Value("cafeFranchise").(*models.CafeFranchise)

	fname, lname := utils.FirstLastName(user.Name)
	err = t.transactionTopic.WriteMessage(ctx, utils.WriteJSONObject[messages.TransactionMidtransMessage](&messages.TransactionMidtransMessage{
		UserID:            user.ID,
		FName:             fname,
		Email:             user.Email,
		Phone:             user.PhoneNumber,
		LName:             lname,
		TransactionCode:   transactionHeader.TransactionCode,
		CafeID:            cafe.ID,
		CafeFranchiseName: cafeFranchise.Name,
		CheckoutList:      request.Checkouts,
	}))

	return nil
}
