package transactions_usecases

import (
	"context"
	"fmt"

	"github.com/RandySteven/CafeConnect/be/entities/models"
)

func (t *transactionWorkflow) checkFranchise(ctx context.Context, franchiseID uint64) (*models.CafeFranchise, error) {
	franchise, err := t.cafeFranchiseRepository.FindByID(ctx, franchiseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get franchise: %w", err)
	}

	if franchise.ID == 0 {
		return nil, fmt.Errorf("franchise not found")
	}

	return franchise, nil
}
