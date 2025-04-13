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
	}
)
