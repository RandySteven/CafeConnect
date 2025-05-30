package responses

import "time"

type (
	AddCartResponse struct {
		ID        string    `json:"id"`
		Action    string    `json:"action"`
		CreatedAt time.Time `json:"created_at"`
	}

	CafeCartItems struct {
		ProductID    uint64     `json:"product_id"`
		ProductName  string     `json:"product_name"`
		ProductImage string     `json:"product_image"`
		ProductPrice uint64     `json:"product_price"`
		Qty          uint64     `json:"qty"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    time.Time  `json:"updated_at"`
		DeletedAt    *time.Time `json:"deleted_at"`
	}

	CheckoutList struct {
		CafeID   uint64           `json:"cafe_id"`
		CafeName string           `json:"cafe_name"`
		Items    []*CafeCartItems `json:"items"`
	}

	ListCartResponse struct {
		UserID       uint64          `json:"user_id"`
		CheckoutList []*CheckoutList `json:"checkout_list"`
	}
)
