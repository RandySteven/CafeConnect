package midtrans_usecases

import (
	"context"
	"log"
	"strconv"

	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/midtrans/midtrans-go"
)

func (m *midtransWorkflow) checkoutList(ctx context.Context, cafeFranchiseName string, checkoutList []*requests.CheckoutList) (result *midtransCheckOut, err error) {
	items := make([]midtrans.ItemDetails, len(checkoutList))
	totalAmount := int64(0)
	for i, item := range checkoutList {
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
			MerchantName: cafeFranchiseName,
		}
		totalAmount += int64(cafeProduct.Price * item.Qty)
	}

	return &midtransCheckOut{
		Items:       items,
		TotalAmount: totalAmount,
	}, nil
}
