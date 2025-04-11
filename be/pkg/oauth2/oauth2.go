package oauth2_client

import (
	"github.com/RandySteven/CafeConnect/be/configs"
	"golang.org/x/oauth2"
)

type (
	Oauth2 interface {
	}

	oauth2Client struct {
		config oauth2.Config
	}
)

func NewOauth2Client(config *configs.Config) (*oauth2Client, error) {
	
	return &oauth2Client{}, nil
}
