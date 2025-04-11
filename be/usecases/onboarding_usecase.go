package usecases

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
)

type onboardingUsecase struct {
	userRepo    repository_interfaces.UserRepository
	pointRepo   repository_interfaces.PointRepository
	transaction repository_interfaces.Transaction
}

func (o *onboardingUsecase) RegisterUser(ctx context.Context, request *requests.RegisterUserRequest) (result *responses.RegisterUserResponse, customErr *apperror.CustomError) {
	var (
		_ = &models.User{}
		_ = &models.Point{}
	)

	customErr = o.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
		//insert user
		_, err := o.userRepo.Save(ctx, nil)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to create user`, err)
		}

		//insert point
		return
	})
	if customErr != nil {
		return nil, customErr
	}
	return
}

func (o *onboardingUsecase) LoginUser(ctx context.Context, request *requests.LoginUserRequest) (result *responses.LoginUserResponse, customErr *apperror.CustomError) {
	return
}

var _ usecase_interfaces.OnboardingUsecase = &onboardingUsecase{}

func newOnboardingUsecase(
	userRepo repository_interfaces.UserRepository,
	pointRepo repository_interfaces.PointRepository,
	transaction repository_interfaces.Transaction,
) *onboardingUsecase {
	return &onboardingUsecase{
		userRepo:    userRepo,
		pointRepo:   pointRepo,
		transaction: transaction,
	}
}
