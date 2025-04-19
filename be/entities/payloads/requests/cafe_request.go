package requests

import (
	"io"
)

type (
	RegisterCafeRequest struct {
		//Name    string `json:"cafe"`
		//Address struct {
		//	Address   string  `json:"address"`
		//	Latitude  float32 `json:"latitude"`
		//	Longitude float32 `json:"longitude"`
		//} `json:"address"`
		//LogoURL     string `json:"logo_url"`
		//CafeType    string `json:"cafe_type"`
		//WorkingHour struct {
		//	OpenHour  time.Time `json:"open_hour"`
		//	CloseHour time.Time `json:"close_hour"`
		//} `json:"working_hour"`
		//PhotoURLs []string `json:"photo_urls"`
		Name       string      `form:"name"`
		Address    string      `form:"address"`
		Latitude   float64     `form:"latitude"`
		Longitude  float64     `form:"longitude"`
		LogoFile   io.Reader   `form:"logo_file"`
		OpenHour   string      `form:"open_hour"`
		CloseHour  string      `form:"close_hour"`
		PhotoFiles []io.Reader `form:"photo_files[]"`
	}

	AddCafeFranchiseRequest struct {
		FranchiseID uint64      `form:"franchise_id"`
		Address     string      `form:"address"`
		Latitude    float64     `form:"latitude"`
		Longitude   float64     `form:"longitude"`
		OpenHour    string      `form:"open_hour"`
		CloseHour   string      `form:"close_hour"`
		PhotoFiles  []io.Reader `form:"photo_files[]"`
	}

	GetCafeListRequest struct {
		Point struct {
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
		} `json:"point"`
		Radius uint64 `json:"radius"`
	}
)
