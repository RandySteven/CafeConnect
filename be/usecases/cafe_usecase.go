package usecases

import (
	"context"
	"errors"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	storage_client "github.com/RandySteven/CafeConnect/be/pkg/storage"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"mime/multipart"
	"sync"
	"time"
)

type cafeUsecase struct {
	transaction   repository_interfaces.Transaction
	cafeRepo      repository_interfaces.CafeRepository
	franchiseRepo repository_interfaces.CafeFranchiseRepository
	addressRepo   repository_interfaces.AddressRepository
	googleStorage storage_client.GoogleStorage
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
	)

	franchise.LogoURL, err = c.googleStorage.UploadFile(ctx, enums.CafesStorage+`logos/`, request.LogoFile, ctx.Value(enums.FileHeader).(*multipart.FileHeader), 40, 40)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `there is issue while upload logo`, err)
	}

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
		}()

		//4. insert address data
		go func() {
			defer wg.Done()
			address, err = c.addressRepo.Save(ctx, address)
			if err != nil {
				customErrCh <- apperror.NewCustomError(apperror.ErrInternalServer, `failed to insert address data`, err)
				return
			}
		}()

		go func() {
			wg.Wait()
			close(customErrCh)
		}()

		//5. insert cafe data
		select {
		case customErr = <-customErrCh:
			return customErr
		default:
			cafe.CafeFranchiseID = franchise.ID
			cafe.AddressID = address.ID
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
	return
}

var _ usecase_interfaces.CafeUsecase = &cafeUsecase{}

func newCafeUsecase(
	cafeRepo repository_interfaces.CafeRepository,
	franchiseRepo repository_interfaces.CafeFranchiseRepository,
	addressRepo repository_interfaces.AddressRepository,
	transaction repository_interfaces.Transaction,
	googleStorage storage_client.GoogleStorage,
	cache cache_interfaces.CafeCache) *cafeUsecase {
	return &cafeUsecase{
		cafeRepo:      cafeRepo,
		franchiseRepo: franchiseRepo,
		addressRepo:   addressRepo,
		transaction:   transaction,
		googleStorage: googleStorage,
		cache:         cache,
	}
}
