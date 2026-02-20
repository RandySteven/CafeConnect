package auto_transfer_usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/enums"
	"github.com/RandySteven/CafeConnect/be/utils"
)

func (t *autoTransferWorkflow) saveTransactionHeader(ctx context.Context, executionData *TransferExecutionData) (*TransferExecutionData, error) {
	transactionHeader := &models.TransactionHeader{
		UserID:          executionData.User.ID,
		TransactionCode: utils.GenerateCode(24),
		CafeID:          executionData.Request.CafeID,
		Status:          enums.TransactionPENDING.String(),
		TransactionAt:   time.Now(),
	}
	transactionHeader, err := t.transactionHeaderRepository.Save(ctx, transactionHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to save transaction header: %w", err)
	}

	executionData.TransactionHeader = transactionHeader
	return executionData, nil
}
