package auto_transfer_usecases

import (
	"context"
	"fmt"
	"time"
)

func (t *autoTransferWorkflow) stockDeduction(ctx context.Context, state *TransferState) (*TransferState, error) {
	for _, checkout := range state.Request.Checkouts {
		cafeProduct, err := t.cafeProductRepository.FindByID(ctx, checkout.CafeProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get cafe product: %w", err)
		}
		if cafeProduct.Stock < checkout.Qty {
			return nil, fmt.Errorf("insufficient stock for product %d", checkout.CafeProductID)
		}
		cafeProduct.Stock -= checkout.Qty
		cafeProduct.UpdatedAt = time.Now()
		_, err = t.cafeProductRepository.Update(ctx, cafeProduct)
		if err != nil {
			return nil, fmt.Errorf("failed to update cafe product: %w", err)
		}
	}
	return state, nil
}
