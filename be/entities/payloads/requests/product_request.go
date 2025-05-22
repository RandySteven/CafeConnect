package requests

import "io"

type (
	AddProductRequest struct {
		Name              string    `form:"name"`
		Photo             io.Reader `form:"photo"`
		FranchiseID       uint64    `form:"franchise_id"`
		ProductCategoryID uint64    `form:"product_category_id"`
		Price             uint64    `form:"price"`
	}

	GetProductListByCafeIDRequest struct {
		CafeID []uint64 `json:"cafe_ids"`
	}
)
