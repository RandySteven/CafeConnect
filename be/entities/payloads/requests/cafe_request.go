package requests

import "time"

type (
	RegisterCafeRequest struct {
		Name    string `json:"cafe"`
		Address struct {
			Address   string  `json:"address"`
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
		} `json:"address"`
		LogoURL     string `json:"logo_url"`
		CafeType    string `json:"cafe_type"`
		WorkingHour struct {
			OpenHour  time.Time `json:"open_hour"`
			CloseHour time.Time `json:"close_hour"`
		} `json:"working_hour"`
		PhotoURLs []string `json:"photo_urls"`
	}

	GetCafeListRequest struct {
		Point struct {
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
		} `json:"point"`
		Radius uint64 `json:"radius"`
	}
)
