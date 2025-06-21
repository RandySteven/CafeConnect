package responses

import "time"

type (
	RegisterCafeResponse struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}

	ListCafeResponse struct {
		ID        uint64 `json:"id"`
		Name      string `json:"name"`
		LogoURL   string `json:"logo_url"`
		Status    string `json:"status"`
		OpenHour  string `json:"open_hour"`
		CloseHour string `json:"close_hour"`
		Address   struct {
			Address   string  `json:"address"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"address"`
	}

	FranchiseListResponse struct {
		ID        uint64     `json:"id"`
		Name      string     `json:"name"`
		LogoURL   string     `json:"logo_url"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}

	FranchiseDetailResponse struct {
		ID            uint64     `json:"id"`
		Name          string     `json:"name"`
		LogoURL       string     `json:"logo_url"`
		NumbOfOutlets uint64     `json:"numb_of_outlets"`
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     time.Time  `json:"updated_at"`
		DeletedAt     *time.Time `json:"deleted_at"`
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
		Status    string     `json:"status"`
		PhotoURLs []string   `json:"photo_urls"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}
)
