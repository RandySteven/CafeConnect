package auto_transfer_usecases

import (
	"context"
	"fmt"
	"time"
)

func (t *autoTransferWorkflow) stockDeduction(ctx context.Context, executionData *TransferExecutionData) (*TransferExecutionData, error) {
	for _, checkout := range executionData.Request.Checkouts {
		cafeProduct, err := t.cafeProductRepository.FindByID(ctx, checkout.CafeProductID)
		if err != nil {
			// Something already deducted — branch to compensation
			executionData.StockDeductionFailed = true
			executionData.NextActivity = autoTransferRestoreStockActivity
			return executionData, nil
		}
		if cafeProduct.Stock < checkout.Qty {
			// Insufficient stock — branch to compensation to undo prior deductions
			executionData.StockDeductionFailed = true
			executionData.NextActivity = autoTransferRestoreStockActivity
			return executionData, nil
			// return nil, fmt.Errorf("insufficient stock for product %d", checkout.CafeProductID)
		}

		cafeProduct.Stock -= checkout.Qty
		cafeProduct.UpdatedAt = time.Now()
		_, err = t.cafeProductRepository.Update(ctx, cafeProduct)
		if err != nil {
			executionData.StockDeductionFailed = true
			executionData.NextActivity = autoTransferRestoreStockActivity
			return executionData, nil
		}

		// Track the successful deduction so restoreStock can undo it
		executionData.DeductedProducts = append(executionData.DeductedProducts, &DeductedProduct{
			CafeProductID: checkout.CafeProductID,
			Qty:           checkout.Qty,
		})
	}
	return executionData, nil
}

func (t *autoTransferWorkflow) restoreStock(ctx context.Context, executionData *TransferExecutionData) (*TransferExecutionData, error) {
	for _, dp := range executionData.DeductedProducts {
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
	executionData.DeductedProducts = nil
	return executionData, nil
}
