package midtrans_usecases

import (
	"context"
	"log"
	"strconv"

	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/midtrans/midtrans-go"
)

func (m *midtransWorkflow) checkoutList(ctx context.Context, executionData *MidtransExecutionData) (*MidtransExecutionData, error) {
	items := make([]midtrans.ItemDetails, len(executionData.Message.CheckoutList))
	totalAmount := int64(0)
	for i, item := range executionData.Message.CheckoutList {
		cafeProduct, err := m.cafeProductRepository.FindByID(ctx, item.CafeProductID)
		if err != nil {
			log.Println("failed to find cafe product:", err)
			return nil, err
		}

		product, err := m.productRepository.FindByID(ctx, cafeProduct.ProductID)
		if err != nil {
			log.Println("failed to find product:", err)
			return nil, err
		}

		items[i] = midtrans.ItemDetails{
			ID:           strconv.FormatUint(item.CafeProductID, 10),
			Name:         product.Name,
			Qty:          int32(item.Qty),
			Price:        int64(cafeProduct.Price),
			MerchantName: executionData.Message.CafeFranchiseName,
		}
		totalAmount += int64(cafeProduct.Price * item.Qty)
	}

	executionData.Items = items
	executionData.TotalAmount = totalAmount
	executionData.MidtransRequest = &midtrans_client.MidtransRequest{
		FName:           executionData.Message.FName,
		LName:           executionData.Message.LName,
		Email:           executionData.Message.Email,
		Phone:           executionData.Message.Phone,
		TransactionCode: executionData.Message.TransactionCode,
		GrossAmt:        totalAmount,
		Items:           items,
	}
	return executionData, nil
}
