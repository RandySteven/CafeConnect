package usecases

import (
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/RandySteven/CafeConnect/be/repositories"
)

type Usecases struct {
	OnboardingUsecase usecase_interfaces.OnboardingUsecase
}

func NewUsecases(repo *repositories.Repositories) *Usecases {
	return &Usecases{}
}
