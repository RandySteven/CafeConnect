package requests

type (
	AddCafeRequest struct {
	}

	GetCafeListRequest struct {
		Point struct {
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
		} `json:"point"`
		Radius uint64 `json:"radius"`
	}
)
