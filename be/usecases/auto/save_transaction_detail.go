package auto_transfer_usecases

import (
	"context"
	"fmt"

	"github.com/RandySteven/CafeConnect/be/entities/models"
)

func (t *autoTransferWorkflow) saveTransactionDetail(ctx context.Context, executionData *TransferExecutionData) (*TransferExecutionData, error) {
	for _, checkout := range executionData.Request.Checkouts {
		transactionDetail := &models.TransactionDetail{
			TransactionID: executionData.TransactionHeader.ID,
			CafeProductID: checkout.CafeProductID,
			Qty:           checkout.Qty,
		}
		_, err := t.transactionDetailRepository.Save(ctx, transactionDetail)
		if err != nil {
			return nil, fmt.Errorf("failed to save transaction detail: %w", err)
		}
	}
	return executionData, nil
}
