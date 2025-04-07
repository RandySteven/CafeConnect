package usecase_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
)

type OnboardingUsecase interface {
	RegisterUser(ctx context.Context, request *requests.RegisterUserRequest)
	LoginUser(ctx context.Context, request *requests.LoginUserRequest)
}
