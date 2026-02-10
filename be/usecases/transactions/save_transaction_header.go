package transactions_usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/enums"
	"github.com/RandySteven/CafeConnect/be/utils"
)

func (t *transactionWorkflow) saveTransactionHeader(ctx context.Context, userID uint64, request *requests.CreateTransactionRequest) (*models.TransactionHeader, error) {
	transactionHeader := &models.TransactionHeader{
		UserID:          userID,
		TransactionCode: utils.GenerateCode(24),
		CafeID:          request.CafeID,
		Status:          enums.TransactionPENDING.String(),
		TransactionAt:   time.Now(),
	}
	transactionHeader, err := t.transactionHeaderRepository.Save(ctx, transactionHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to save transaction header: %w", err)
	}

	return transactionHeader, nil
}
