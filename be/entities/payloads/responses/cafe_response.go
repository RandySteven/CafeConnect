package responses

type (
	AddCafeResponse struct {
	}

	ListCafeResponse struct {
		ID          uint64 `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		LogoURL     string `json:"logo_url"`
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
