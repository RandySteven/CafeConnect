package requests

type (
	AddAddressRequest struct {
		Address   string  `json:"address"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		IsDefault bool    `json:"is_default"`
		OwnerType string  `json:"owner_type"`
	}
)
