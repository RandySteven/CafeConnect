package auto_transfer_usecases

import (
	"context"
	"fmt"
)

func (t *autoTransferWorkflow) checkUser(ctx context.Context, executionData *TransferExecutionData) (*TransferExecutionData, error) {
	user, err := t.userRepository.FindByID(ctx, executionData.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	executionData.User = user
	return executionData, nil
}
