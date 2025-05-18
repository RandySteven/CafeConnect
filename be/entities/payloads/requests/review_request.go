package requests

type (
	AddReviewRequest struct {
		CafeID  uint64  `json:"cafe_id"`
		Score   float64 `json:"score"`
		Comment string  `json:"comment"`
	}

	GetCafeReviewRequest struct {
		CafeID uint64 `json:"cafe_id"`
	}
)
