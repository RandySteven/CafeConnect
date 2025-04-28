package requests

import "io"

type (
	AddProductRequest struct {
		Name              string    `form:"name"`
		Photo             io.Reader `form:"photo"`
		ProductType       string    `form:"product_type"`
		FranchiseID       uint64    `form:"franchise_id"`
		ProductCategoryID uint64    `form:"product_category_id"`
		Price             uint64    `form:"price"`
	}
)
