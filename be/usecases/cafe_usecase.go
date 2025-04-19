package usecases

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
)

type cafeUsecase struct {
	transaction repository_interfaces.Transaction
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

func newCafeUsecase(transaction repository_interfaces.Transaction) *cafeUsecase {
	return &cafeUsecase{
		transaction: transaction,
	}
}
