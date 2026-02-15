package auto_transfer_usecases

import (
	"context"
	"fmt"
)

func (t *autoTransferWorkflow) checkFranchise(ctx context.Context, state *TransferState) (*TransferState, error) {
	franchise, err := t.cafeFranchiseRepository.FindByID(ctx, state.Cafe.CafeFranchiseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get franchise: %w", err)
	}

	if franchise.ID == 0 {
		return nil, fmt.Errorf("franchise not found")
	}

	state.Franchise = franchise
	return state, nil
}
