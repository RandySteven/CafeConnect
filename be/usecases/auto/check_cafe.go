package auto_transfer_usecases

import (
	"context"
	"fmt"
)

func (t *autoTransferWorkflow) checkCafe(ctx context.Context, executionData *ExecutionData) (*ExecutionData, error) {
	cafe, err := t.cafeRepository.FindByID(ctx, executionData.Request.CafeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cafe: %w", err)
	}

	if cafe.ID == 0 {
		return nil, fmt.Errorf("cafe not found")
	}

	executionData.Cafe = cafe
	executionData.SetActivity(autoTransferCheckFranchiseActivity)
	return executionData, nil
}
