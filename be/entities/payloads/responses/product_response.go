package responses

import (
	"time"
)

type (
	AddProductResponse struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}

	DetailProductResponse struct {
		ID              uint64 `json:"id"`
		Name            string `json:"name"`
		Photo           string `json:"photo"`
		Price           uint64 `json:"price"`
		ProductCategory *struct {
			ID       uint64 `json:"id"`
			Category string `json:"category"`
		} `json:"product_category"`
		Stock     uint64    `json:"stock"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		DeletedAt time.Time `json:"deleted_at"`
	}

	ListProductResponse struct {
		ID        uint64    `json:"id"`
		Name      string    `json:"name"`
		Photo     string    `json:"photo"`
		Price     uint64    `json:"price"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		DeletedAt time.Time `json:"deleted_at"`
	}
)
