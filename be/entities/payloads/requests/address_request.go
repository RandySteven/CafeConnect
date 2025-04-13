package requests

type (
	AddAddress struct {
		OwnerID   uint64  `json:"owner_id"`
		Address   string  `json:"address"`
		Longitude float32 `json:"longitude"`
		Latitude  float32 `json:"latitude"`
		IsDefault bool    `json:"is_default"`
		OwnerType string  `json:"owner_type"`
	}
)
