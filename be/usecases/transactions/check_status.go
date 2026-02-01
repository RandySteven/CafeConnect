package transactions_usecases

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
)

func (t *transactionWorkflow) checkStatus(ctx context.Context) (err error) {
	transactionHeader := ctx.Value("transactionHeader").(*models.TransactionHeader)

	transactionHeader, err = t.transactionHeaderRepository.FindByTransactionCode(ctx, transactionHeader.TransactionCode)
	if err != nil {
		return err
	}
	var midtransTransaction *models.MidtransTransaction

	if transactionHeader.Status == enums.TransactionSUCCESS.String() {
		midtransTransaction, err = t.midtransTransactionRepository.FindByTransactionCode(ctx, transactionHeader.TransactionCode)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	} else {
		return nil
	}

	result := &responses.TransactionReceiptResponse{
		ID:              transactionHeader.ID,
		TransactionCode: transactionHeader.TransactionCode,
		Status:          transactionHeader.Status,
		TransactionAt:   transactionHeader.TransactionAt.Local(),
		MidtransResponse: &midtrans_client.MidtransResponse{
			Token:       midtransTransaction.Token,
			RedirectURL: midtransTransaction.RedirectURL,
		},
	}
	ctx = context.WithValue(ctx, "result", result)

	return nil
}
