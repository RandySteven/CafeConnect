package usecases

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/google/uuid"
	"sync"
)

type onboardingUsecase struct {
	userRepo     repository_interfaces.UserRepository
	pointRepo    repository_interfaces.PointRepository
	referralRepo repository_interfaces.ReferralRepository
	transaction  repository_interfaces.Transaction
}

func (o *onboardingUsecase) RegisterUser(ctx context.Context, request *requests.RegisterUserRequest) (result *responses.RegisterUserResponse, customErr *apperror.CustomError) {
	var (
		user = &models.User{
			Name:           fmt.Sprintf("%s %s", request.FirstName, request.LastName),
			Username:       request.Username,
			Email:          request.Email,
			ProfilePicture: request.ProfilePicture,
			PhoneNumber:    request.PhoneNumber,
		}
		point = &models.Point{
			Point: 0,
		}
		referral = &models.Referral{
			Code: uuid.NewString(),
		}
		err         error
		wg          *sync.WaitGroup
		customErrCh = make(chan *apperror.CustomError)
	)

	customErr = o.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
		//insert user
		user, err = o.userRepo.Save(ctx, user)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to create user`, err)
		}

		//check referral
		if request.Referral != nil {

		}

		wg.Add(3)
		//parallelysm progress between point and referral
		//insert point
		go func() {
			defer wg.Done()
			point.UserID = user.ID
			point, err = o.pointRepo.Save(ctx, point)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to create point`, err)
				return
			}
		}()

		//insert referral
		go func() {
			defer wg.Done()
			referral.UserID = user.ID
			referral, err = o.referralRepo.Save(ctx, referral)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to create point`, err)
				return
			}
		}()

		//insert address
		go func() {
			defer wg.Done()

		}()
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
