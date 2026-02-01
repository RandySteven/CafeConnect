package transactions_usecases

import (
	"context"
	"errors"

	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
)

func (t *transactionWorkflow) checkFranchise(ctx context.Context, request *requests.CreateTransactionRequest) (err error) {
	cafe := ctx.Value("cafe").(*models.Cafe)

	franchise, err := t.cafeFranchiseRepository.FindByID(ctx, cafe.ID)
	if err != nil {
		return err
	}

	if franchise.ID == 0 {
		return errors.New("franchise not found")
	}

	return nil
}
