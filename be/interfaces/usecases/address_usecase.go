package usecase_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type AddressUsecase interface {
	AddUserAddress(ctx context.Context, request *requests.AddAddressRequest) (result *responses.AddAddressResponse, customErr *apperror.CustomError)
	GetUserAddresses(ctx context.Context) (result []*responses.GetUserAddressResponse, customErr *apperror.CustomError)
}
