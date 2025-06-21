package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	aws_client "github.com/RandySteven/CafeConnect/be/pkg/aws"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"mime/multipart"
	"strings"
	"sync"
	"time"
)

type cafeUsecase struct {
	transaction   repository_interfaces.Transaction
	cafeRepo      repository_interfaces.CafeRepository
	franchiseRepo repository_interfaces.CafeFranchiseRepository
	addressRepo   repository_interfaces.AddressRepository
	aws           aws_client.AWS
	cache         cache_interfaces.CafeCache
}

func (c *cafeUsecase) GetListCafeFranchises(ctx context.Context) (result []*responses.FranchiseListResponse, customErr *apperror.CustomError) {
	result, err := c.cache.GetFranchiseListCache(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed in redis`, err)
	}
	if result != nil {
		return result, nil
	}
	franchises, err := c.franchiseRepo.FindAll(ctx, 0, 0)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get franchises`, err)
	}
	for _, franchise := range franchises {
		result = append(result, &responses.FranchiseListResponse{
			ID:        franchise.ID,
			Name:      franchise.Name,
			LogoURL:   franchise.LogoURL,
			CreatedAt: franchise.CreatedAt,
			UpdatedAt: franchise.UpdatedAt,
			DeletedAt: franchise.DeletedAt,
		})
	}
	c.cache.SetFranchiseListCache(ctx, result)
	return result, nil
}

func (c *cafeUsecase) RegisterCafeAndFranchise(ctx context.Context, request *requests.RegisterCafeAndFranchiseRequest) (result *responses.RegisterCafeResponse, customErr *apperror.CustomError) {
	//1. proceed the request
	var (
		cafe      = &models.Cafe{}
		franchise = &models.CafeFranchise{
			Name: request.Name,
		}
		address = &models.Address{
			Address:   request.Address,
			Longitude: request.Longitude,
			Latitude:  request.Latitude,
		}
		err           error
		wg            sync.WaitGroup
		numbOfWorkers = 2
		customErrCh   = make(chan *apperror.CustomError)
		addressIdCh   = make(chan uint64)
		franchiseIdCh = make(chan uint64)
		resultPaths   = []string{}
	)

	logoResult, err := c.aws.UploadImageFile(ctx, request.LogoFile, enums.CafesStorage+`logos/`, ctx.Value(enums.FileHeader).(*multipart.FileHeader), 40, 40)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `there is issue while upload logo`, err)
	}
	franchise.LogoURL = *logoResult
	for _, file := range request.PhotoFiles {
		resultPath, err := c.aws.UploadImageFile(ctx, file, fmt.Sprintf("%s%s/", enums.CafesStorage, utils.CafeNameToSnakeCase(request.Name)), ctx.Value(enums.FileHeader).(*multipart.FileHeader), 1920, 1080)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `there is issue while upload photo`, err)
		}
		resultPaths = append(resultPaths, *resultPath)
	}

	photoUrls := utils.Join(resultPaths, ";")

	//2. create transaction
	customErr = c.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
		wg.Add(numbOfWorkers)
		//3. insert cafe franchise data
		go func() {
			defer wg.Done()
			franchise, err = c.franchiseRepo.Save(ctx, franchise)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to insert franchise data`, err)
				return
			}
			franchiseIdCh <- franchise.ID
		}()

		//4. insert address data
		go func() {
			defer wg.Done()
			address, err = c.addressRepo.Save(ctx, address)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to insert address data`, err)
				return
			}
			addressIdCh <- address.ID
		}()

		go func() {
			wg.Wait()
			close(customErrCh)
			close(addressIdCh)
			close(franchiseIdCh)
		}()

		//5. insert cafe data
		select {
		case customErr = <-customErrCh:

			wg.Add(1)
			go func() {
				defer wg.Done()
				//for _, resultPath := range resultPaths {
				//	_ = c.googleStorage.DeleteFile(ctx, resultPath)
				//}
			}()
			wg.Wait()

			return customErr
		default:
			cafe.CafeFranchiseID = <-franchiseIdCh
			cafe.AddressID = <-addressIdCh
			cafe.PhotoURLs = photoUrls
			cafe.CafeType = request.CafeType
			cafe.CloseHour = request.CloseHour
			cafe.OpenHour = request.OpenHour
			cafe, err = c.cafeRepo.Save(ctx, cafe)
			if err != nil {
				return apperror.NewCustomError(apperror.ErrInternalServer, `failed to create cafe`, err)
			}
			return nil
		}
	})
	if customErr != nil {
		return nil, customErr
	}

	return &responses.RegisterCafeResponse{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
	}, nil
}

func (c *cafeUsecase) GetListOfCafeBasedOnRadius(ctx context.Context, request *requests.GetCafeListRequest) (result []*responses.ListCafeResponse, customErr *apperror.CustomError) {
	key := fmt.Sprintf(enums.ListCafeRadiusKey, request.Point.Longitude, request.Point.Latitude, fmt.Sprintf("%d", request.Radius))

	result, err := c.cache.GetCafeRadiusListCache(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cache`, err)
	}

	if result != nil {
		return result, nil
	}

	targetAddress := request.Point

	addresses, err := c.addressRepo.FindAddressBasedOnRadius(ctx, targetAddress.Longitude, targetAddress.Latitude, request.Radius)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get address`, err)
	}

	for _, address := range addresses {
		cafe, err := c.cafeRepo.FindByAddressId(ctx, address.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe`, err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			continue
		}

		cafeFranchise, err := c.franchiseRepo.FindByID(ctx, cafe.CafeFranchiseID)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get franchise`, err)
		}

		result = append(result, &responses.ListCafeResponse{
			ID:        cafe.ID,
			Name:      cafeFranchise.Name,
			Status:    utils.GetCafeOpenCloseStatus(cafe.OpenHour, cafe.CloseHour),
			LogoURL:   utils.ImageStorage(cafeFranchise.LogoURL),
			OpenHour:  cafe.OpenHour[0:5],
			CloseHour: cafe.CloseHour[0:5],
			Address: struct {
				Address   string  `json:"address"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			}{Address: address.Address, Latitude: address.Latitude, Longitude: address.Longitude},
		})
	}
	_ = c.cache.SetCafeRadiusListCache(ctx, key, result)
	return result, nil
}

func (c *cafeUsecase) GetCafeDetail(ctx context.Context, id uint64) (result *responses.DetailCafeResponse, customErr *apperror.CustomError) {
	var (
		wg          sync.WaitGroup
		customErrCh = make(chan *apperror.CustomError)
		franchiseCh = make(chan *models.CafeFranchise)
		addressCh   = make(chan *models.Address)
		err         error
		cafe        *models.Cafe
		franchise   *models.CafeFranchise
		address     *models.Address
	)

	result, err = c.cache.GetCafeDetail(ctx, fmt.Sprintf("%d", id))
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get redis detail`, err)
	}
	if result != nil {
		return result, nil
	}

	cafe, err = c.cafeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe`, err)
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		franchise, err = c.franchiseRepo.FindByID(ctx, cafe.CafeFranchiseID)
		if err != nil {
			customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed franchise`, err)
			return
		}
		franchiseCh <- franchise
	}()

	go func() {
		defer wg.Done()
		address, err = c.addressRepo.FindByID(ctx, cafe.AddressID)
		if err != nil {
			customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to get address`, err)
			return
		}
		addressCh <- address
	}()

	go func() {
		wg.Wait()
		close(customErrCh)
		close(addressCh)
		close(franchiseCh)
	}()

	select {
	case customErr = <-customErrCh:
		return nil, customErr
	default:
		address = <-addressCh
		franchise = <-franchiseCh
		result = &responses.DetailCafeResponse{
			ID:      cafe.ID,
			Name:    franchise.Name,
			LogoURL: utils.ImageStorage(franchise.LogoURL),
			Status:  utils.GetCafeOpenCloseStatus(cafe.OpenHour, cafe.CloseHour),
			Address: struct {
				Address   string  `json:"address"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			}{Address: address.Address, Latitude: address.Latitude, Longitude: address.Longitude},
			PhotoURLs: strings.Split(cafe.PhotoURLs, ";"),
			CreatedAt: cafe.CreatedAt,
			UpdatedAt: cafe.UpdatedAt,
			DeletedAt: cafe.DeletedAt,
		}
		c.cache.SetCafeDetail(ctx, fmt.Sprintf("%d", result.ID), result)
		return result, nil
	}
}

func (c *cafeUsecase) AddCafeOutlet(ctx context.Context, request *requests.AddCafeOutletRequest) (result *responses.RegisterCafeResponse, customErr *apperror.CustomError) {
	var (
		cafe      = &models.Cafe{}
		franchise = &models.CafeFranchise{}
		address   = &models.Address{
			Address:   request.Address,
			Longitude: request.Longitude,
			Latitude:  request.Latitude,
		}
		err         error
		resultPaths = []string{}
	)

	franchise, err = c.franchiseRepo.FindByID(ctx, request.FranchiseID)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get franchise`, err)
	}

	if len(request.PhotoFiles) != 0 {
		for _, file := range request.PhotoFiles {
			resultPath, err := c.aws.UploadImageFile(ctx, file, fmt.Sprintf("%s%s/", enums.CafesStorage, utils.CafeNameToSnakeCase(franchise.Name)), ctx.Value(enums.FileHeader).(*multipart.FileHeader), 1920, 1080)
			if err != nil {
				return nil, apperror.NewCustomError(apperror.ErrInternalServer, `there is issue while upload photo`, err)
			}
			resultPaths = append(resultPaths, *resultPath)
		}
	}

	photoUrls := utils.Join(resultPaths, ";")

	if customErr = c.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {

		//1. insert address
		address, err = c.addressRepo.Save(ctx, address)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to add address`, err)
		}

		//2. insert cafe
		cafe = &models.Cafe{
			AddressID:       address.ID,
			CafeFranchiseID: franchise.ID,
			CafeType:        "",
			PhotoURLs:       photoUrls,
			OpenHour:        request.OpenHour,
			CloseHour:       request.CloseHour,
		}
		cafe, err = c.cafeRepo.Save(ctx, cafe)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to add cafe`, err)
		}

		return nil
	}); customErr != nil {
		return nil, customErr
	}

	return &responses.RegisterCafeResponse{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
	}, nil
}

var _ usecase_interfaces.CafeUsecase = &cafeUsecase{}

func newCafeUsecase(
	cafeRepo repository_interfaces.CafeRepository,
	franchiseRepo repository_interfaces.CafeFranchiseRepository,
	addressRepo repository_interfaces.AddressRepository,
	transaction repository_interfaces.Transaction,
	aws aws_client.AWS,
	cache cache_interfaces.CafeCache) *cafeUsecase {
	return &cafeUsecase{
		cafeRepo:      cafeRepo,
		franchiseRepo: franchiseRepo,
		addressRepo:   addressRepo,
		transaction:   transaction,
		aws:           aws,
		cache:         cache,
	}
}
