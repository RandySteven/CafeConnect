package auto_transfer_usecases

import (
	"context"
	"fmt"

	"github.com/RandySteven/CafeConnect/be/entities/models"
)

func (t *autoTransferWorkflow) checkCafe(ctx context.Context, cafeID uint64) (*models.Cafe, error) {
	cafe, err := t.cafeRepository.FindByID(ctx, cafeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cafe: %w", err)
	}

	if cafe.ID == 0 {
		return nil, fmt.Errorf("cafe not found")
	}

	return cafe, nil
}
