package usecase_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type CafeUsecase interface {
	RegisterCafeAndFranchise(ctx context.Context, request *requests.RegisterCafeAndFranchiseRequest) (result *responses.RegisterCafeResponse, customErr *apperror.CustomError)
	GetListOfCafeBasedOnRadius(ctx context.Context, request *requests.GetCafeListRequest) (result []*responses.ListCafeResponse, customErr *apperror.CustomError)
	GetListCafeFranchises(ctx context.Context) (result []*responses.FranchiseListResponse, customErr *apperror.CustomError)
	GetCafeDetail(ctx context.Context, id uint64) (result *responses.DetailCafeResponse, customErr *apperror.CustomError)
	AddCafeOutlet(ctx context.Context, request *requests.AddCafeOutletRequest) (result *responses.RegisterCafeResponse, customErr *apperror.CustomError)
}
