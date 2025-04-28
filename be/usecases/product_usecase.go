package usecases

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
)

type productUsecase struct {
}

func (p *productUsecase) GetProductByCafe(ctx context.Context, cafeId uint64) (result []*responses.ListProductResponse, customErr *apperror.CustomError) {
	//TODO implement me
	panic("implement me")
}

func (p *productUsecase) GetProductDetail(ctx context.Context, id uint64) (result *responses.DetailProductResponse, customErr *apperror.CustomError) {
	//TODO implement me
	panic("implement me")
}

var _ usecase_interfaces.ProductUsecase = &productUsecase{}

func newProductUsecase() *productUsecase {
	return &productUsecase{}
}
