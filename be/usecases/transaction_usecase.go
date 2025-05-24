package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type transactionUsecase struct {
	transactionHeaderRepository repository_interfaces.TransactionHeaderRepository
	transactionDetailRepository repository_interfaces.TransactionDetailRepository
	cartRepository              repository_interfaces.CartRepository
	userRepository              repository_interfaces.UserRepository
	cafeRepository              repository_interfaces.CafeRepository
	productRepository           repository_interfaces.ProductRepository
	cafeProductRepository       repository_interfaces.CafeProductRepository
	transaction                 repository_interfaces.Transaction
	midtrans                    midtrans_client.Midtrans
	cache                       cache_interfaces.TransactionCache
}

func (t *transactionUsecase) CreateTransaction(ctx context.Context) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError) {
	userId := ctx.Value(enums.UserID).(uint64)
	var (
		carts             []*models.Cart
		err               error
		transactionHeader *models.TransactionHeader
		transactionDetail *models.TransactionDetail
		cafeId            uint64 = 0
		ids               []uint64
		mu                sync.Mutex
		cafeProduct       = &models.CafeProduct{}
	)
	//1. check cart is not empty
	carts, err = t.cartRepository.FindByUserID(ctx, userId)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get user cart`, err)
	}
	for _, cart := range carts {
		ids = append(ids, cart.ID)
	}

	if customErr = t.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
		//pre cond (1, get cafe id)
		cafeId, err = t.cafeProductRepository.FindCafeIDByCafeProductIDs(ctx, ids)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe id`, err)
		}

		//2. create transaction header
		transactionHeader = &models.TransactionHeader{
			UserID:          userId,
			CafeID:          cafeId,
			TransactionCode: utils.GenerateCode(24),
			Status:          `SUCCESS`,
			TransactionAt:   time.Now(),
		}

		transactionHeader, err = t.transactionHeaderRepository.Save(ctx, transactionHeader)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to insert transaction header`, err)
		}

		for _, cart := range carts {

			//3. reduce qty of cafe_product
			mu.Lock()
			cafeProduct, err = t.cafeProductRepository.FindByID(ctx, cart.CafeProductID)
			if err != nil {
				return apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe product`, err)
			}
			if cafeProduct.Stock < cart.Qty {
				return apperror.NewCustomError(apperror.ErrBadRequest, `insufficient product stock`, fmt.Errorf(`product stock less than qty`))
			}
			cafeProduct.Stock -= cart.Qty
			if cafeProduct.Stock == 0 {
				cafeProduct.Status = `NOT_AVAILABLE`
			}
			cafeProduct.UpdatedAt = time.Now()
			cafeProduct, err = t.cafeProductRepository.Update(ctx, cafeProduct)
			if err != nil {
				return apperror.NewCustomError(apperror.ErrInternalServer, `failed to update cafe product`, err)
			}
			mu.Unlock()

			//4. create transaction detail
			transactionDetail = &models.TransactionDetail{
				TransactionID: transactionHeader.ID,
				CafeProductID: cafeProduct.ID,
				Qty:           cart.Qty,
			}
			transactionDetail, err = t.transactionDetailRepository.Save(ctx, transactionDetail)
			if err != nil {
				return apperror.NewCustomError(apperror.ErrInternalServer, `failed to insert transaction detail`, err)
			}
		}

		//5. delete cart
		err = t.cartRepository.DeleteByUserID(ctx, userId)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to delete cart`, err)
		}
		return nil
	}); customErr != nil {
		return nil, customErr
	}

	result = &responses.TransactionReceiptResponse{
		ID:              transactionHeader.ID,
		TransactionCode: transactionHeader.TransactionCode,
		TransactionAt:   transactionHeader.TransactionAt,
		Status:          transactionHeader.Status,
	}

	return
}

func (t *transactionUsecase) GetUserTransactions(ctx context.Context) (result []*responses.TransactionListResponse, customErr *apperror.CustomError) {
	//TODO implement me
	panic("implement me")
}

func (t *transactionUsecase) GetTransactionByCode(ctx context.Context, transactionCode string) (result *responses.TransactionDetailResponse, customErr *apperror.CustomError) {
	var (
		err                error
		transactionHeader  *models.TransactionHeader
		transactionDetail  *models.TransactionDetail
		transactionDetails []*models.TransactionDetail
		cafeProduct        *models.CafeProduct
		product            *models.Product
		item               *responses.TransactionDetailItem
	)
	result = &responses.TransactionDetailResponse{}

	result, err = t.cache.Get(ctx, transactionCode)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get transaction detail by redis`, err)
	}

	if result != nil {
		return result, nil
	}

	err = t.transactionHeaderRepository.CreateIndex(ctx)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to create index`, err)
	}

	transactionHeader, err = t.transactionHeaderRepository.FindByTransactionCode(ctx, transactionCode)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get transaction header`, err)
	}

	result.ID = transactionHeader.ID
	result.TransactionCode = transactionHeader.TransactionCode
	result.TransactionTime = transactionHeader.TransactionAt
	result.Status = transactionHeader.Status

	transactionDetails, err = t.transactionDetailRepository.FindByTransactionId(ctx, transactionHeader.ID)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get transaction details`, err)
	}

	for _, transactionDetail = range transactionDetails {
		cafeProduct, err = t.cafeProductRepository.FindByID(ctx, transactionDetail.CafeProductID)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe product`, err)
		}

		product, err = t.productRepository.FindByID(ctx, cafeProduct.ProductID)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get product`, err)
		}

		item = &responses.TransactionDetailItem{
			ID:       cafeProduct.ID,
			Name:     product.Name,
			Price:    cafeProduct.Price,
			ImageURL: product.PhotoURL,
			Qty:      transactionDetail.Qty,
		}

		result.Items = append(result.Items, item)
	}

	err = t.transactionHeaderRepository.DropIndex(ctx)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to drop index`, err)
	}

	_ = t.cache.Set(ctx, transactionCode, result)

	return result, nil
}

var _ usecase_interfaces.TransactionUsecase = &transactionUsecase{}

func newTransactionUsecase(
	transactionHeaderRepository repository_interfaces.TransactionHeaderRepository,
	transactionDetailRepository repository_interfaces.TransactionDetailRepository,
	cartRepository repository_interfaces.CartRepository,
	userRepository repository_interfaces.UserRepository,
	cafeRepository repository_interfaces.CafeRepository,
	productRepository repository_interfaces.ProductRepository,
	cafeProductRepository repository_interfaces.CafeProductRepository,
	transaction repository_interfaces.Transaction,
	cache cache_interfaces.TransactionCache,
	midtrans midtrans_client.Midtrans) *transactionUsecase {
	return &transactionUsecase{
		transactionHeaderRepository: transactionHeaderRepository,
		transactionDetailRepository: transactionDetailRepository,
		cartRepository:              cartRepository,
		userRepository:              userRepository,
		cafeRepository:              cafeRepository,
		productRepository:           productRepository,
		cafeProductRepository:       cafeProductRepository,
		transaction:                 transaction,
		cache:                       cache,
		midtrans:                    midtrans,
	}
}
