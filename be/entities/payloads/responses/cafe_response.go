package responses

import "time"

type (
	AddCafeResponse struct {
	}

	ListCafeResponse struct {
		ID        uint64    `json:"id"`
		Name      string    `json:"name"`
		LogoURL   string    `json:"logo_url"`
		Status    string    `json:"status"`
		OpenHour  time.Time `json:"open_hour"`
		CloseHour time.Time `json:"close_hour"`
	}

	FranchiseListResponse struct {
		ID        uint64     `json:"id"`
		Name      string     `json:"name"`
		LogoURL   string     `json:"logo_url"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}

	DetailCafeResponse struct {
		ID      uint64 `json:"id"`
		Name    string `json:"name"`
		LogoURL string `json:"logo_url"`
		Address struct {
			Address   string  `json:"address"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"address"`
		PhotoURLs []string `json:"photo_urls"`
	}
)
