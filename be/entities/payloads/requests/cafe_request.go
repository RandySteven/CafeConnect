package requests

import (
	"io"
)

type (
	RegisterCafeAndFranchiseRequest struct {
		Name       string      `form:"name"`
		Address    string      `form:"address"`
		Latitude   float64     `form:"latitude"`
		Longitude  float64     `form:"longitude"`
		LogoFile   io.Reader   `form:"logo_file"`
		OpenHour   string      `form:"open_hour"`
		CloseHour  string      `form:"close_hour"`
		PhotoFiles []io.Reader `form:"photo_files[]"`
		CafeType   string      `form:"cafe_type"`
	}

	AddCafeOutletRequest struct {
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
			Longitude float64 `json:"longitude"`
			Latitude  float64 `json:"latitude"`
		} `json:"point"`
		Radius uint64 `json:"radius"`
	}
)
