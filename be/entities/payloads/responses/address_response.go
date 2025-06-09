package responses

import "time"

type (
	AddAddressResponse struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}

	GetUserAddressResponse struct {
		ID        uint64     `json:"id"`
		Address   string     `json:"address"`
		Latitude  float64    `json:"latitude"`
		Longitude float64    `json:"longitude"`
		IsDefault bool       `json:"is_default"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}
)
