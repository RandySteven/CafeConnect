package usecases

import (
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	storage_client "github.com/RandySteven/CafeConnect/be/pkg/storage"
	"github.com/RandySteven/CafeConnect/be/repositories"
)

type Usecases struct {
	OnboardingUsecase usecase_interfaces.OnboardingUsecase
}

func NewUsecases(repo *repositories.Repositories,
	googleStorage storage_client.GoogleStorage) *Usecases {
	return &Usecases{}
}
