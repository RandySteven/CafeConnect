package usecase_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type ProductUsecase interface {
	GetProductByCafe(ctx context.Context, cafeId uint64) (result []*responses.ListProductResponse, customErr *apperror.CustomError)
	GetProductDetail(ctx context.Context, id uint64) (result *responses.DetailProductResponse, customErr *apperror.CustomError)
}
