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
	var (
		user = &models.User{
			Name:        fmt.Sprintf("%s %s", request.FirstName, request.LastName),
			Username:    request.Username,
			Email:       request.Email,
			PhoneNumber: request.PhoneNumber,
			DoB:         utils.ConvertDateString(request.DoB),
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
		err         error
		wg          *sync.WaitGroup
		customErrCh = make(chan *apperror.CustomError)
	)

	//photo uploader
	fileHeader := ctx.Value(enums.FileHeader).(*multipart.FileHeader)
	o.googleStorage.UploadFile(ctx, "", request.ProfilePicture, fileHeader, 40, 40)

	customErr = o.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
		//insert user
		user, err = o.userRepo.Save(ctx, user)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to create user`, err)
		}

		//check referral
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

		wg.Add(3)
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
			referral.Code = uuid.NewString()
			referral, err = o.referralRepo.Save(ctx, referral)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to create referall`, err)
				return
			}
		}()

		//insert address
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
	transaction repository_interfaces.Transaction,
	googleStorage storage_client.GoogleStorage,
) *onboardingUsecase {
	return &onboardingUsecase{
		userRepo:      userRepo,
		pointRepo:     pointRepo,
		transaction:   transaction,
		googleStorage: googleStorage,
	}
}
