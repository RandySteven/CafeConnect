package apis

import (
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	"github.com/RandySteven/CafeConnect/be/usecases"
)

type APIs struct {
	DevApi        api_interfaces.DevApi
	OnboardingApi api_interfaces.OnboardingApi
}

func NewAPIs(usecases *usecases.Usecases) *APIs {
	return &APIs{
		DevApi:        newDevApi(),
		OnboardingApi: newOnboardingApi(usecases.OnboardingUsecase),
	}
}
