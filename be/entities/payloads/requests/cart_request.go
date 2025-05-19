package requests

type (
	AddToCartRequest struct {
		CafeProductID uint64 `json:"cafe_product_id"`
		Qty           uint64 `json:"qty"`
	}
)
