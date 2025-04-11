package apis

import (
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	"github.com/RandySteven/CafeConnect/be/usecases"
)

type APIs struct {
	DevApi api_interfaces.DevApi
}

func NewAPIs(usecases usecases.Usecases) *APIs {
	return &APIs{
		DevApi: newDevApi(),
	}
}
