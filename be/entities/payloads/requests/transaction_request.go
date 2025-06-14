package requests

type (
	CheckoutList struct {
		CafeProductID uint64 `json:"cafe_product_id"`
		Qty           uint64 `json:"qty"`
	}

	CreateTransactionRequest struct {
		CafeID    uint64          `json:"cafe_id"`
		Checkouts []*CheckoutList `json:"checkouts"`
	}

	ReceiptRequest struct {
		TransactionCode string `json:"transaction_code"`
	}
)
