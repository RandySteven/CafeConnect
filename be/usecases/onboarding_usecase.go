package usecases

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	storage_client "github.com/RandySteven/CafeConnect/be/pkg/storage"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"sync"
	"time"
)

type onboardingUsecase struct {
	userRepo        repository_interfaces.UserRepository
	pointRepo       repository_interfaces.PointRepository
	referralRepo    repository_interfaces.ReferralRepository
	addressRepo     repository_interfaces.AddressRepository
	addressUserRepo repository_interfaces.AddressUserRepository
	transaction     repository_interfaces.Transaction
	googleStorage   storage_client.GoogleStorage
}

func (o *onboardingUsecase) RegisterUser(ctx context.Context, request *requests.RegisterUserRequest) (result *responses.RegisterUserResponse, customErr *apperror.CustomError) {
	var err error
	fileHeader := ctx.Value(enums.FileHeader).(*multipart.FileHeader)
	resultPath, err := o.googleStorage.UploadFile(ctx, enums.UsersStorage, request.ProfilePicture, fileHeader, 40, 40)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrBadRequest, `failed to upload image `, err)
	}
	timeDoB, err := utils.ConvertDateString(request.DoB)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrBadRequest, `failed to convert dob `, err)
	}
	log.Println("time dob ", *timeDoB)

	var (
		password = utils.HashPassword(request.Password)

		user = &models.User{
			Name:           fmt.Sprintf("%s %s", request.FirstName, request.LastName),
			Username:       request.Username,
			Email:          request.Email,
			PhoneNumber:    request.PhoneNumber,
			DoB:            time.DateOnly,
			Password:       password,
			ProfilePicture: resultPath,
		}
		point = &models.Point{
			Point: 0,
		}
		referral = &models.Referral{}
		address  = &models.Address{
			Address:   request.Address,
			Latitude:  request.Latitude,
			Longitude: request.Longitude,
		}
		wg          sync.WaitGroup
		customErrCh = make(chan *apperror.CustomError)
	)

	customErr = o.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
		numbOfWorkers := 3
		user, err = o.userRepo.Save(ctx, user)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to create user`, err)
		}

		if request.ReferralCode != "" {
			referral, err = o.referralRepo.FindByCode(ctx, request.ReferralCode)
			if err != nil {
				return apperror.NewCustomError(apperror.ErrInternalServer, `failed to get referral`, err)
			}
			point.Point += 100
			referral.NumbOfUsage += 1
			referral.UpdatedAt = time.Now()
			_, err = o.referralRepo.Update(ctx, referral)
			if err != nil {
				return apperror.NewCustomError(apperror.ErrInternalServer, `failed to update referral`, err)
			}
		}

		wg.Add(numbOfWorkers)
		go func() {
			defer wg.Done()
			referral.UserID = user.ID
			referral.Code = utils.RandomString(16)
			referral.ExpiredTime = time.Now().Add(8 * 24 * time.Hour)
			referral, err = o.referralRepo.Save(ctx, referral)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to create referall`, err)
				return
			}
		}()

		go func() {
			defer wg.Done()
			point.UserID = user.ID
			point, err = o.pointRepo.Save(ctx, point)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to create point`, err)
				return
			}
		}()

		go func() {
			defer wg.Done()
			address, err = o.addressRepo.Save(ctx, address)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to insert address`, err)
				return
			}
			addressUser := &models.AddressUser{
				AddressID: address.ID,
				UserID:    user.ID,
			}
			_, err = o.addressUserRepo.Save(ctx, addressUser)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to insert address user`, err)
				return
			}
		}()

		go func() {
			wg.Wait()
			close(customErrCh)
		}()

		select {
		case customErr = <-customErrCh:
			return customErr
		}
	})
	if customErr != nil {
		return nil, customErr
	}
	return &responses.RegisterUserResponse{
		ID:           uuid.NewString(),
		Email:        request.Email,
		RegisterTime: time.Now(),
	}, nil
}

func (o *onboardingUsecase) LoginUser(ctx context.Context, request *requests.LoginUserRequest) (result *responses.LoginUserResponse, customErr *apperror.CustomError) {
	return
}

var _ usecase_interfaces.OnboardingUsecase = &onboardingUsecase{}

func newOnboardingUsecase(
	userRepo repository_interfaces.UserRepository,
	pointRepo repository_interfaces.PointRepository,
	addressRepo repository_interfaces.AddressRepository,
	addressUserRepo repository_interfaces.AddressUserRepository,
	referralRepo repository_interfaces.ReferralRepository,
	transaction repository_interfaces.Transaction,
	googleStorage storage_client.GoogleStorage,
) *onboardingUsecase {
	return &onboardingUsecase{
		userRepo:        userRepo,
		pointRepo:       pointRepo,
		addressRepo:     addressRepo,
		addressUserRepo: addressUserRepo,
		referralRepo:    referralRepo,
		transaction:     transaction,
		googleStorage:   googleStorage,
	}
}
