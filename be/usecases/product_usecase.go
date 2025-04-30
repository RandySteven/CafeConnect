package usecases

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
)

type productUsecase struct {
	cafeRepo                  repository_interfaces.CafeRepository
	cafeFranchiseRepo         repository_interfaces.CafeFranchiseRepository
	cafeProductRepo           repository_interfaces.CafeProductRepository
	productRepo               repository_interfaces.ProductRepository
	productCategoryRepository repository_interfaces.ProductCategoryRepository
	transaction               repository_interfaces.Transaction
}

func (p *productUsecase) GetProductByCafe(ctx context.Context, cafeId uint64) (result []*responses.ListProductResponse, customErr *apperror.CustomError) {
	panic("implement me")
}

func (p *productUsecase) GetProductDetail(ctx context.Context, id uint64) (result *responses.DetailProductResponse, customErr *apperror.CustomError) {
	panic("implement me")
}

var _ usecase_interfaces.ProductUsecase = &productUsecase{}

func newProductUsecase(
	cafeRepo repository_interfaces.CafeRepository,
	cafeFranchiseRepo repository_interfaces.CafeFranchiseRepository,
	cafeProductRepo repository_interfaces.CafeProductRepository,
	productRepo repository_interfaces.ProductRepository,
	productCategoryRepository repository_interfaces.ProductCategoryRepository,
	transaction repository_interfaces.Transaction,
) *productUsecase {
	return &productUsecase{
		cafeRepo:                  cafeRepo,
		cafeFranchiseRepo:         cafeFranchiseRepo,
		cafeProductRepo:           cafeProductRepo,
		productRepo:               productRepo,
		productCategoryRepository: productCategoryRepository,
		transaction:               transaction,
	}
}
