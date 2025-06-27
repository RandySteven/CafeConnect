package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	aws_client "github.com/RandySteven/CafeConnect/be/pkg/aws"
	jwt_client "github.com/RandySteven/CafeConnect/be/pkg/jwt"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"sync"
	"time"
)

type onboardingUsecase struct {
	userRepo        repository_interfaces.UserRepository
	pointRepo       repository_interfaces.PointRepository
	referralRepo    repository_interfaces.ReferralRepository
	addressRepo     repository_interfaces.AddressRepository
	addressUserRepo repository_interfaces.AddressUserRepository
	verifyTokenRepo repository_interfaces.VerifyTokenRepository
	transaction     repository_interfaces.Transaction
	onboardingCache cache_interfaces.OnboardingCache
	onboardingTopic topics_interfaces.OnboardingTopic
	aws             aws_client.AWS
}

func (o *onboardingUsecase) RegisterUser(ctx context.Context, request *requests.RegisterUserRequest) (result *responses.RegisterUserResponse, customErr *apperror.CustomError) {
	var (
		err        error
		dummyImg   = os.Getenv(`DEFAULT_DUMMY_IMG`)
		resultPath = &dummyImg
	)

	if request.ProfilePicture != nil {
		fileHeader := ctx.Value(enums.FileHeader).(*multipart.FileHeader)
		resultPath, err = o.aws.UploadImageFile(ctx, request.ProfilePicture, enums.UsersStorage, fileHeader, 0, 0)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrBadRequest, `failed to upload image `, err)
		}
	}
	timeDoB, err := utils.ConvertDateString(request.DoB)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrBadRequest, `failed to convert dob `, err)
	}

	var (
		password = utils.HashPassword(request.Password)

		user = &models.User{
			Name:           fmt.Sprintf("%s %s", request.FirstName, request.LastName),
			Username:       request.Username,
			Email:          request.Email,
			PhoneNumber:    request.PhoneNumber,
			DoB:            timeDoB.String(),
			Password:       password,
			ProfilePicture: *resultPath,
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
			//_ = o.googleStorage.DeleteFile(ctx, resultPath)
			return customErr
		}
	})
	if customErr != nil {
		return nil, customErr
	}
	_ = o.onboardingTopic.WriteMessage(ctx, utils.WriteJSONObject[messages.VerifyTokenMessage](&messages.VerifyTokenMessage{
		Token:  ctx.Value(enums.RequestID).(string),
		UserID: user.ID,
	}))
	return &responses.RegisterUserResponse{
		ID:           uuid.NewString(),
		Email:        request.Email,
		RegisterTime: time.Now(),
	}, nil
}

func (o *onboardingUsecase) LoginUser(ctx context.Context, request *requests.LoginUserRequest) (result *responses.LoginUserResponse, customErr *apperror.CustomError) {
	user, err := o.userRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewCustomError(apperror.ErrNotFound, `failed to login consumers not found`, err)
		}
		log.Println(err)
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to connect db`, err)
	}
	isPassExists := utils.ComparePassword(request.Password, user.Password)
	if !isPassExists {
		return nil, apperror.NewCustomError(apperror.ErrNotFound, `invalid credentials`, err)
	}

	accessToken, refreshToken := jwt_client.GenerateTokens(user, nil)

	result = &responses.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		LoginTime:    time.Now(),
	}
	return result, nil
}

func (o *onboardingUsecase) GetOnboardUser(ctx context.Context) (result *responses.OnboardUserResponse, customErr *apperror.CustomError) {
	id := ctx.Value(enums.UserID).(uint64)
	result = &responses.OnboardUserResponse{}
	var (
		user                   = &models.User{}
		point                  = &models.Point{}
		addressUsers           = []*models.AddressUser{}
		err                    error
		wg                     sync.WaitGroup
		customErrCh            = make(chan *apperror.CustomError)
		addressUsersResponse   = []*responses.OnboardUserAddress{}
		addressUsersResponseCh = make(chan []*responses.OnboardUserAddress)
		userCh                 = make(chan *models.User)
		pointCh                = make(chan *models.Point)
	)
	result, err = o.onboardingCache.Get(ctx, strconv.Itoa(int(id)))
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed on progress the cache`, err)
	}
	if result != nil {
		log.Println("INI DAPAT DARI REDIS LOH")
		return result, nil
	}

	numbOfWorkers := 3
	wg.Add(numbOfWorkers)
	go func() {
		defer wg.Done()
		user, err = o.userRepo.FindByID(ctx, id)
		if err != nil {
			customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to get onboard user`, err)
			return
		}
		userCh <- user
	}()

	go func() {
		defer wg.Done()
		point, err = o.pointRepo.FindByUserID(ctx, id)
		if err != nil {
			customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to get point user `, err)
			return
		}
		pointCh <- point
	}()

	go func() {
		defer wg.Done()
		addressUsers, err = o.addressUserRepo.FindByUserID(ctx, id)
		if err != nil {
			customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to get address users`, err)
			return
		}
		for _, addressUser := range addressUsers {
			address, err := o.addressRepo.FindByID(ctx, addressUser.AddressID)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to get address `, err)
				return
			}
			addressUsersResponse = append(addressUsersResponse, &responses.OnboardUserAddress{
				ID:        address.ID,
				Address:   address.Address,
				Latitude:  address.Latitude,
				Longitude: address.Longitude,
				IsDefault: addressUser.IsDefault,
			})
		}
		addressUsersResponseCh <- addressUsersResponse
	}()

	go func() {
		wg.Wait()
		close(customErrCh)
		close(userCh)
		close(pointCh)
		close(addressUsersResponseCh)
	}()

	select {
	case customErr = <-customErrCh:
		return nil, customErr
	default:
		user = <-userCh
		point = <-pointCh
		addressUsersResponse = <-addressUsersResponseCh
		result = &responses.OnboardUserResponse{
			ID:             user.ID,
			Name:           user.Name,
			Username:       user.Username,
			ProfilePicture: user.ProfilePicture,
			Email:          user.Email,
			Point:          point.Point,
			Addresses:      addressUsersResponse,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			DeletedAt:      user.DeletedAt,
		}
		o.onboardingCache.Set(ctx, strconv.Itoa(int(id)), result)
		return result, nil
	}
}

func (o *onboardingUsecase) VerifyUser(ctx context.Context, tokenID string) (customErr *apperror.CustomError) {
	verifyToken, err := o.verifyTokenRepo.FindByToken(ctx, tokenID)
	if err != nil {
		return apperror.NewCustomError(apperror.ErrNotFound, `not found verify token`, err)
	}

	if time.Now().After(verifyToken.ExpiredTime) {
		return apperror.NewCustomError(apperror.ErrForbidden, `token already expired`, fmt.Errorf(`expired token`))
	}

	verifyToken.IsClicked = true
	verifyToken.UpdatedAt = time.Now()
	_, err = o.verifyTokenRepo.Update(ctx, verifyToken)
	if err != nil {
		return apperror.NewCustomError(apperror.ErrInternalServer, `failed to update token`, err)
	}

	return nil
}

var _ usecase_interfaces.OnboardingUsecase = &onboardingUsecase{}

func newOnboardingUsecase(
	userRepo repository_interfaces.UserRepository,
	pointRepo repository_interfaces.PointRepository,
	addressRepo repository_interfaces.AddressRepository,
	addressUserRepo repository_interfaces.AddressUserRepository,
	referralRepo repository_interfaces.ReferralRepository,
	verifyTokenRepo repository_interfaces.VerifyTokenRepository,
	transaction repository_interfaces.Transaction,
	onboardingCache cache_interfaces.OnboardingCache,
	onboardingTopic topics_interfaces.OnboardingTopic,
	aws aws_client.AWS,
) *onboardingUsecase {
	return &onboardingUsecase{
		userRepo:        userRepo,
		pointRepo:       pointRepo,
		addressRepo:     addressRepo,
		addressUserRepo: addressUserRepo,
		verifyTokenRepo: verifyTokenRepo,
		referralRepo:    referralRepo,
		transaction:     transaction,
		onboardingCache: onboardingCache,
		onboardingTopic: onboardingTopic,
		aws:             aws,
	}
}
