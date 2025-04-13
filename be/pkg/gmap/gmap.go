package gmap_client

import (
	"github.com/RandySteven/CafeConnect/be/configs"
	"googlemaps.github.io/maps"
)

type (
	gmapClient struct {
		client *maps.Client
	}
)

func NewGMapClient(config *configs.Config) (*gmapClient, error) {
	c, err := maps.NewClient(maps.WithAPIKey(``))
	if err != nil {
		return nil, err
	}
	return &gmapClient{
		client: c,
	}, nil
}
