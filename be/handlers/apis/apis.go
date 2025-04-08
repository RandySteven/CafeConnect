package apis

import "github.com/RandySteven/CafeConnect/be/usecases"

type APIs struct {
}

func NewAPIs(usecases usecases.Usecases) *APIs {
	return &APIs{}
}
