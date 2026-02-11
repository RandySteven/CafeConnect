package auto_transfer_usecases

import (
	"context"
	"fmt"

	"github.com/RandySteven/CafeConnect/be/entities/models"
)

func (t *autoTransferWorkflow) checkUser(ctx context.Context, userID uint64) (*models.User, error) {
	user, err := t.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
