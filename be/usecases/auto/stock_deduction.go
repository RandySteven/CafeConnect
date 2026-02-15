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
			// Something already deducted — branch to compensation
			state.StockDeductionFailed = true
			state.NextActivity = autoTransferRestoreStockActivity
			return state, nil
		}
		if cafeProduct.Stock < checkout.Qty {
			// Insufficient stock — branch to compensation to undo prior deductions
			state.StockDeductionFailed = true
			state.NextActivity = autoTransferRestoreStockActivity
			return state, nil
		}

		cafeProduct.Stock -= checkout.Qty
		cafeProduct.UpdatedAt = time.Now()
		_, err = t.cafeProductRepository.Update(ctx, cafeProduct)
		if err != nil {
			state.StockDeductionFailed = true
			state.NextActivity = autoTransferRestoreStockActivity
			return state, nil
		}

		// Track the successful deduction so restoreStock can undo it
		state.DeductedProducts = append(state.DeductedProducts, &DeductedProduct{
			CafeProductID: checkout.CafeProductID,
			Qty:           checkout.Qty,
		})
	}
	return state, nil
}

func (t *autoTransferWorkflow) restoreStock(ctx context.Context, state *TransferState) (*TransferState, error) {
	for _, dp := range state.DeductedProducts {
		cafeProduct, err := t.cafeProductRepository.FindByID(ctx, dp.CafeProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to find product %d for stock restoration: %w", dp.CafeProductID, err)
		}

		cafeProduct.Stock += dp.Qty
		cafeProduct.UpdatedAt = time.Now()
		_, err = t.cafeProductRepository.Update(ctx, cafeProduct)
		if err != nil {
			return nil, fmt.Errorf("failed to restore stock for product %d: %w", dp.CafeProductID, err)
		}
	}

	// Clear after restoration
	state.DeductedProducts = nil
	return state, nil
}
