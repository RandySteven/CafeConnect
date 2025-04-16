package responses

import "time"

type (
	RegisterUserResponse struct {
		ID           string    `json:"id"`
		Email        string    `json:"email"`
		RegisterTime time.Time `json:"register_time"`
	}

	LoginUserResponse struct {
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		LoginTime    time.Time `json:"login_time"`
	}
)
