package transactions_usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
)

func (t *transactionWorkflow) transactionCheckStatus(ctx context.Context, transactionCode string) (*responses.TransactionReceiptResponse, error) {
	transactionHeader, err := t.transactionHeaderRepository.FindByTransactionCode(ctx, transactionCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction header: %w", err)
	}

	if transactionHeader.Status != enums.TransactionSUCCESS.String() {
		return nil, fmt.Errorf("transaction is not successful, status: %s", transactionHeader.Status)
	}

	midtransTransaction, err := t.midtransTransactionRepository.FindByTransactionCode(ctx, transactionHeader.TransactionCode)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get midtrans transaction: %w", err)
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

	return result, nil
}
