package usecases

import usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"

type Usecases struct {
	OnboardingUsecase usecase_interfaces.OnboardingUsecase
}

func NewUsecases() *Usecases {
	return &Usecases{}
}
