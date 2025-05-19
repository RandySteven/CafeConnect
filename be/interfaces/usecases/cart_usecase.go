package usecase_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type CartUsecase interface {
	AddToCart(ctx context.Context, request *requests.AddToCartRequest) (result *responses.AddCartResponse, customErr *apperror.CustomError)
	GetUserCart(ctx context.Context) (result *responses.ListCartResponse, customErr *apperror.CustomError)
}
