package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/redis/go-redis/v9"
	"log"
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
	pub                         kafka_client.Publisher
	midtrans                    midtrans_client.Midtrans
	cache                       cache_interfaces.TransactionCache
}

func (t *transactionUsecase) CheckoutTransactionV1(ctx context.Context) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError) {
	//TODO implement me
	panic("implement me")
}

func (t *transactionUsecase) CheckoutTransactionV2(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError) {
	//process read the request
	var (
		userId            = ctx.Value(enums.UserID).(uint64)
		transactionHeader = &models.TransactionHeader{
			UserID:          userId,
			TransactionCode: utils.GenerateCode(24),
			CafeID:          request.CafeID,
			Status:          enums.TransactionPENDING.String(),
			TransactionAt:   time.Now(),
		}
		transactionDetail *models.TransactionDetail
		err               error
	)

	if customErr = t.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
		transactionHeader, err = t.transactionHeaderRepository.Save(ctx, transactionHeader)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to create header transaction`, err)
		}

		for _, item := range request.Checkouts {
			transactionDetail = &models.TransactionDetail{
				TransactionID: transactionHeader.ID,
				CafeProductID: item.CafeProductID,
				Qty:           item.Qty,
			}

			transactionDetail, err = t.transactionDetailRepository.Save(ctx, transactionDetail)
			if err != nil {
				return apperror.NewCustomError(apperror.ErrInternalServer, `failed to create detail transaction`, err)
			}

		}
		err = t.cartRepository.DeleteByUserID(ctx, userId)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to delete cart`, err)
		}

		transactionHeader, err = t.transactionHeaderRepository.FindByTransactionCode(ctx, transactionHeader.TransactionCode)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to get transaction header`, err)
		}

		transactionHeader.TransactionCode = enums.TransactionSUCCESS.String()
		transactionHeader.UpdatedAt = time.Now()
		transactionHeader, err = t.transactionHeaderRepository.Update(ctx, transactionHeader)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to update trans status`, err)
		}
		return nil
	}); customErr != nil {
		return nil, customErr
	}

	return &responses.TransactionReceiptResponse{
		ID:              transactionHeader.ID,
		TransactionCode: transactionHeader.TransactionCode,
		Status:          transactionHeader.Status,
		TransactionAt:   transactionHeader.TransactionAt,
	}, nil
}

func (t *transactionUsecase) CreateTransactionV1(ctx context.Context) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError) {
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

	if len(carts) == 0 {
		return nil, apperror.NewCustomError(apperror.ErrBadRequest, `cart still empty`, fmt.Errorf(`cart still empty`))
	}

	for _, cart := range carts {
		ids = append(ids, cart.CafeProductID)
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

func (t *transactionUsecase) CreateTransactionV2(ctx context.Context) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError) {
	userId := ctx.Value(enums.UserID).(uint64)

	var (
		carts             []*models.Cart
		err               error
		transactionHeader *models.TransactionHeader
		cafeId            uint64 = 0
		ids               []uint64
	)

	user, err := t.userRepository.FindByID(ctx, userId)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get user`, err)
	}

	carts, err = t.cartRepository.FindByUserID(ctx, userId)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get user cart`, err)
	}
	if len(carts) == 0 {
		return nil, apperror.NewCustomError(apperror.ErrBadRequest, `cart still empty`, fmt.Errorf(`cart still empty`))
	}

	for _, cart := range carts {
		ids = append(ids, cart.CafeProductID)
	}

	transactionHeader = &models.TransactionHeader{
		UserID:          userId,
		CafeID:          cafeId,
		TransactionCode: utils.GenerateCode(24),
		Status:          enums.TransactionPENDING.String(),
		TransactionAt:   time.Now(),
	}

	transactionHeader, err = t.transactionHeaderRepository.Save(ctx, transactionHeader)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to create transaction header`, err)
	}

	firstName, lastName := utils.FirstLastName(user.Name)
	midtransRequest := &messages.TransactionMidtransMessage{
		UserID:      userId,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: user.PhoneNumber,
	}

	t.pub.WriteMessage(ctx, enums.TransactionTopic, `midtrans-request`, utils.WriteJSONObject[messages.TransactionMidtransMessage](midtransRequest))

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

	result, err = t.cache.Get(ctx, transactionCode)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get transaction detail by redis`, err)
	}

	if result != nil {
		return result, nil
	}

	//err = t.transactionHeaderRepository.CreateIndex(ctx)
	//if err != nil {
	//	return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to create index`, err)
	//}

	transactionHeader, err = t.transactionHeaderRepository.FindByTransactionCode(ctx, transactionCode)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get transaction header`, err)
	}
	log.Println(`transaction header`, transactionHeader)
	result = &responses.TransactionDetailResponse{
		ID:              transactionHeader.ID,
		TransactionCode: transactionHeader.TransactionCode,
		TransactionTime: transactionHeader.TransactionAt,
		Status:          transactionHeader.Status,
	}

	if result.Status == enums.TransactionPENDING.String() {
		return result, nil
	}

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

	//err = t.transactionHeaderRepository.DropIndex(ctx)
	//if err != nil {
	//	return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to drop index`, err)
	//}

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
