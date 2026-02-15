package auto_transfer_usecases

import (
	"context"
	"fmt"
)

func (t *autoTransferWorkflow) checkUser(ctx context.Context, state *TransferState) (*TransferState, error) {
	user, err := t.userRepository.FindByID(ctx, state.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	state.User = user
	return state, nil
}
