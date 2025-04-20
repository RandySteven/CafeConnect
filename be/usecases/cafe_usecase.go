package usecases

import (
	"context"
	"errors"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/redis/go-redis/v9"
)

type cafeUsecase struct {
	transaction   repository_interfaces.Transaction
	cafeRepo      repository_interfaces.CafeRepository
	franchiseRepo repository_interfaces.CafeFranchiseRepository
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

func (c *cafeUsecase) RegisterCafe(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (c *cafeUsecase) GetListOfCafeBasedOnRadius(ctx context.Context, request *requests.GetCafeListRequest) (result []*responses.ListCafeResponse, customErr *apperror.CustomError) {
	//TODO implement me
	panic("implement me")
}

var _ usecase_interfaces.CafeUsecase = &cafeUsecase{}

func newCafeUsecase(
	cafeRepo repository_interfaces.CafeRepository,
	franchiseRepo repository_interfaces.CafeFranchiseRepository,
	transaction repository_interfaces.Transaction,
	cache cache_interfaces.CafeCache) *cafeUsecase {
	return &cafeUsecase{
		cafeRepo:      cafeRepo,
		franchiseRepo: franchiseRepo,
		transaction:   transaction,
		cache:         cache,
	}
}
