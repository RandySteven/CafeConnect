package auto_transfer_usecases

import (
	"context"
	"fmt"
	"time"
)

func (t *autoTransferWorkflow) updateHeaderStatus(ctx context.Context, transactionCode string, status string) (err error) {
	transactionHeader, err := t.transactionHeaderRepository.FindByTransactionCode(ctx, transactionCode)
	if err != nil {
		return fmt.Errorf("failed to get transaction header: %w", err)
	}
	transactionHeader.Status = status
	transactionHeader.UpdatedAt = time.Now()
	_, err = t.transactionHeaderRepository.Update(ctx, transactionHeader)
	if err != nil {
		return fmt.Errorf("failed to update header status: %w", err)
	}
	return nil
}
