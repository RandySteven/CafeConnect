package auto_transfer_usecases

import (
	"context"
	"fmt"
)

func (t *autoTransferWorkflow) checkFranchise(ctx context.Context, executionData *TransferExecutionData) (*TransferExecutionData, error) {
	franchise, err := t.cafeFranchiseRepository.FindByID(ctx, executionData.Cafe.CafeFranchiseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get franchise: %w", err)
	}

	if franchise.ID == 0 {
		return nil, fmt.Errorf("franchise not found")
	}

	executionData.Franchise = franchise
	return executionData, nil
}
