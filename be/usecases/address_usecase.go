package usecases

import (
	"context"
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
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type addressUsecase struct {
	addressRepo     repository_interfaces.AddressRepository
	addressUserRepo repository_interfaces.AddressUserRepository
	userRepo        repository_interfaces.UserRepository
	addressCache    cache_interfaces.AddressCache
}

func (a *addressUsecase) AddUserAddress(ctx context.Context, request *requests.AddAddressRequest) (result *responses.AddAddressResponse, customErr *apperror.CustomError) {
	userId := ctx.Value(enums.UserID).(uint64)

	address := &models.Address{
		Address:   request.Address,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
	}

	address, err := a.addressRepo.Save(ctx, address)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to save address`, err)
	}

	addressUser := &models.AddressUser{
		AddressID: address.ID,
		UserID:    userId,
		IsDefault: request.IsDefault,
	}

	addressUser, err = a.addressUserRepo.Save(ctx, addressUser)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to save user address`, err)
	}

	return &responses.AddAddressResponse{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
	}, nil
}

func (a *addressUsecase) GetUserAddresses(ctx context.Context) (result []*responses.GetUserAddressResponse, customErr *apperror.CustomError) {
	userId := ctx.Value(enums.UserID).(uint64)
	key := fmt.Sprintf(enums.AddressUserCacheKey, fmt.Sprintf("%d", userId))

	result, err := a.addressCache.GetMultiData(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get data from redis`, err)
	}

	if result != nil {
		return result, nil
	}

	addressUsers, err := a.addressUserRepo.FindByUserID(ctx, userId)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get address user`, err)
	}

	result = make([]*responses.GetUserAddressResponse, len(addressUsers))

	for index, addressUser := range addressUsers {
		address, err := a.addressRepo.FindByID(ctx, addressUser.AddressID)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get address`, err)
		}
		result[index] = &responses.GetUserAddressResponse{
			ID:        address.ID,
			Address:   address.Address,
			Latitude:  address.Latitude,
			Longitude: address.Longitude,
			IsDefault: addressUser.IsDefault,
			CreatedAt: address.CreatedAt,
			UpdatedAt: address.UpdatedAt,
			DeletedAt: address.DeletedAt,
		}
	}

	_ = a.addressCache.SetMultiData(ctx, key, result)

	return result, nil
}

var _ usecase_interfaces.AddressUsecase = &addressUsecase{}

func newAddressUsecase(
	addressRepo repository_interfaces.AddressRepository,
	addressUserRepo repository_interfaces.AddressUserRepository,
	userRepo repository_interfaces.UserRepository,
	addressCache cache_interfaces.AddressCache,
) *addressUsecase {
	return &addressUsecase{
		addressRepo:     addressRepo,
		addressUserRepo: addressUserRepo,
		userRepo:        userRepo,
		addressCache:    addressCache,
	}
}
