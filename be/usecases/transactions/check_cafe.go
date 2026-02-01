package transactions_usecases

import (
	"context"
	"errors"

	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
)

func (t *transactionWorkflow) checkCafe(ctx context.Context, request *requests.CreateTransactionRequest) (err error) {
	cafe, err := t.cafeRepository.FindByID(ctx, request.CafeID)
	if err != nil {
		return err
	}

	if cafe.ID == 0 {
		return errors.New("cafe not found")
	}

	ctx = context.WithValue(ctx, "cafe", cafe)

	return nil
}
