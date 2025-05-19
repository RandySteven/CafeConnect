package responses

import "time"

type (
	AddCartResponse struct {
		ID        string    `json:"id"`
		Action    string    `json:"action"`
		CreatedAt time.Time `json:"created_at"`
	}

	CartItem struct {
		ProductID    uint64     `json:"product_id"`
		ProductName  string     `json:"product_name"`
		ProductPrice uint64     `json:"product_price"`
		Qty          uint64     `json:"qty"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    time.Time  `json:"updated_at"`
		DeletedAt    *time.Time `json:"deleted_at"`
	}

	ListCartResponse struct {
		UserID uint64      `json:"user_id"`
		Items  []*CartItem `json:"items"`
	}
)
