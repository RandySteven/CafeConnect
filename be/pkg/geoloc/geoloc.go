package geoloc_client

import (
	"github.com/RandySteven/CafeConnect/be/configs"
	"net/http"
	"time"
)

type (
	geolocAPI struct {
		client *http.Client
		apiKey string
	}

	Point struct {
		Long float32
		Lat  float32
	}
)

func NewGeolocAPI(config *configs.Config) (*geolocAPI, error) {
	return &geolocAPI{
		client: &http.Client{Timeout: 10 * time.Second},
		apiKey: config.Config.Oauth2.GoogleClientID,
	}, nil
}
