package usecases

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	storage_client "github.com/RandySteven/CafeConnect/be/pkg/storage"
	"github.com/google/uuid"
	"mime/multipart"
	"time"
)

type productUsecase struct {
	cafeRepo                  repository_interfaces.CafeRepository
	cafeFranchiseRepo         repository_interfaces.CafeFranchiseRepository
	cafeProductRepo           repository_interfaces.CafeProductRepository
	productRepo               repository_interfaces.ProductRepository
	productCategoryRepository repository_interfaces.ProductCategoryRepository
	storage                   storage_client.GoogleStorage
	transaction               repository_interfaces.Transaction
}

func (p *productUsecase) AddProduct(ctx context.Context, request *requests.AddProductRequest) (result *responses.AddProductResponse, customErr *apperror.CustomError) {
	var (
		product = &models.Product{
			Name:              request.Name,
			ProductCategoryID: request.ProductCategoryID,
		}
		err   error
		cafes []*models.Cafe
	)

	resultPath, err := p.storage.UploadFile(ctx, enums.ProductsStorage, "", request.Photo, ctx.Value(enums.FileHeader).(*multipart.FileHeader), 40, 40)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to upload product`, err)
	}

	product.PhotoURL = resultPath

	cafes, err = p.cafeRepo.FindByCafeFranchiseId(ctx, request.FranchiseID)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafes`, err)
	}

	customErr = p.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
		product, err = p.productRepo.Save(ctx, product)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to insert product`, err)
		}

		for _, cafe := range cafes {
			_, err := p.cafeProductRepo.Save(ctx, &models.CafeProduct{
				CafeID:    cafe.ID,
				ProductID: product.ID,
				Price:     request.Price,
			})
			if err != nil {
				return apperror.NewCustomError(apperror.ErrInternalServer, `failed to insert cafe product`, err)
			}
		}

		return nil
	})
	if customErr != nil {
		return nil, customErr
	}

	return &responses.AddProductResponse{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
	}, nil
}

func (p *productUsecase) GetProductByCafe(ctx context.Context, cafeId uint64) (result []*responses.ListProductResponse, customErr *apperror.CustomError) {
	var ()

	return
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
	storage storage_client.GoogleStorage,
	transaction repository_interfaces.Transaction,
) *productUsecase {
	return &productUsecase{
		cafeRepo:                  cafeRepo,
		cafeFranchiseRepo:         cafeFranchiseRepo,
		cafeProductRepo:           cafeProductRepo,
		productRepo:               productRepo,
		storage:                   storage,
		productCategoryRepository: productCategoryRepository,
		transaction:               transaction,
	}
}
