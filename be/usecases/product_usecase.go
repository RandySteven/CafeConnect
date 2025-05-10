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
	var (
		cafe         = &models.Cafe{}
		cafeProducts = []*models.CafeProduct{}
		err          error
		product      = &models.Product{}
	)

	cafe, err = p.cafeRepo.FindByID(ctx, cafeId)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe`, err)
	}

	cafeProducts, err = p.cafeProductRepo.FindByCafeID(ctx, cafe.ID)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe products`, err)
	}

	for _, cafeProduct := range cafeProducts {
		product, err = p.productRepo.FindByID(ctx, cafeProduct.ProductID)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get product`, err)
		}

		result = append(result, &responses.ListProductResponse{
			ID:        cafeProduct.ID,
			Name:      product.Name,
			Photo:     product.PhotoURL,
			Price:     cafeProduct.Price,
			CreatedAt: cafeProduct.CreatedAt,
			UpdatedAt: cafeProduct.UpdatedAt,
			DeletedAt: cafeProduct.DeletedAt,
		})
	}

	return result, nil
}

func (p *productUsecase) GetProductDetail(ctx context.Context, id uint64) (result *responses.DetailProductResponse, customErr *apperror.CustomError) {
	var (
		cafeProduct = &models.CafeProduct{}
		product     = &models.Product{}
		err         error
		category    = &models.ProductCategory{}
	)
	cafeProduct, err = p.cafeProductRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe product`, err)
	}

	product, err = p.productRepo.FindByID(ctx, cafeProduct.ProductID)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get product`, err)
	}

	category, err = p.productCategoryRepository.FindByID(ctx, product.ProductCategoryID)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get category`, err)
	}

	result = &responses.DetailProductResponse{
		ID:    cafeProduct.ID,
		Name:  product.Name,
		Photo: product.PhotoURL,
		Price: cafeProduct.Price,
		Stock: cafeProduct.Stock,
		ProductCategory: &struct {
			ID       uint64 `json:"id"`
			Category string `json:"category"`
		}{
			ID:       category.ID,
			Category: category.Category,
		},
		CreatedAt: cafeProduct.CreatedAt,
		UpdatedAt: cafeProduct.UpdatedAt,
		DeletedAt: cafeProduct.DeletedAt,
	}

	return result, nil
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
