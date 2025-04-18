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

	OnboardUserAddress struct {
		ID        uint64  `json:"id"`
		Address   string  `json:"address"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		IsDefault bool    `json:"is_default"`
	}

	OnboardUserResponse struct {
		ID        uint64                `json:"id"`
		Name      string                `json:"name"`
		Username  string                `json:"username"`
		Email     string                `json:"email"`
		Point     uint64                `json:"point"`
		Addresses []*OnboardUserAddress `json:"addresses"`
		CreatedAt time.Time             `json:"created_at"`
		UpdatedAt time.Time             `json:"updated_at"`
		DeletedAt *time.Time            `json:"deleted_at"`
	}
)
