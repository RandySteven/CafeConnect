package requests

import "io"

type (
	AddProductRequest struct {
		Name              string    `form:"name"`
		Photo             io.Reader `form:"photo"`
		ProductCategoryID uint64    `form:"product_category_id"`
		Price             uint64    `form:"price"`
	}
)
