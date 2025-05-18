package responses

import "time"

type (
	AddReviewResponse struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}

	ReviewsResponse struct {
		User struct {
			ID             uint64 `json:"id"`
			Name           string `json:"name"`
			ProfilePicture string `json:"profile_picture"`
		} `json:"user"`
		Score     float64   `json:"score"`
		Comment   string    `json:"comment"`
		CreatedAt time.Time `json:"created_at"`
	}

	GetReviewsResponse struct {
		CafeID   uint64             `json:"cafe_id"`
		AvgScore float64            `json:"avg_score"`
		Reviews  []*ReviewsResponse `json:"reviews"`
	}
)
