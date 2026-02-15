package auto_transfer_usecases

import (
	"context"
	"fmt"
)

func (t *autoTransferWorkflow) checkCafe(ctx context.Context, state *TransferState) (*TransferState, error) {
	cafe, err := t.cafeRepository.FindByID(ctx, state.Request.CafeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cafe: %w", err)
	}

	if cafe.ID == 0 {
		return nil, fmt.Errorf("cafe not found")
	}

	state.Cafe = cafe
	return state, nil
}
