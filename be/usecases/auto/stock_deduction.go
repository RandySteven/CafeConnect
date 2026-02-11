package auto_transfer_usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
)

func (t *autoTransferWorkflow) stockDeduction(ctx context.Context, checkoutList []*requests.CheckoutList) error {
	for _, checkout := range checkoutList {
		cafeProduct, err := t.cafeProductRepository.FindByID(ctx, checkout.CafeProductID)
		if err != nil {
			return fmt.Errorf("failed to get cafe product: %w", err)
		}
		if cafeProduct.Stock < checkout.Qty {
			return fmt.Errorf("insufficient stock: %w", errors.New("insufficient stock"))
		}
		// if index == 2 {
		// 	return fmt.Errorf("failed to stock deduction: %w", errors.New("failed to stock deduction"))
		// }
		cafeProduct.Stock -= checkout.Qty
		cafeProduct.UpdatedAt = time.Now()
		_, err = t.cafeProductRepository.Update(ctx, cafeProduct)
		if err != nil {
			return fmt.Errorf("failed to update cafe product: %w", err)
		}
	}
	return nil
}
