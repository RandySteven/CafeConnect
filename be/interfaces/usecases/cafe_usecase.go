package usecase_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type CafeUsecase interface {
	RegisterCafe(ctx context.Context)
	GetListOfCafeBasedOnRadius(ctx context.Context, request *requests.GetCafeListRequest) (result []*responses.ListCafeResponse, customErr *apperror.CustomError)
	GetListCafeFranchises(ctx context.Context) (result []*responses.FranchiseListResponse, customErr *apperror.CustomError)
}
