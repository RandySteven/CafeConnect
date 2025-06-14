package consumers

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/midtrans/midtrans-go"
	"log"
	"strconv"
	"strings"
	"time"
)

type TransactionConsumer struct {
	consumer                      kafka_client.Consumer
	publisher                     kafka_client.Publisher
	midtrans                      midtrans_client.Midtrans
	userRepository                repository_interfaces.UserRepository
	transactionRepository         repository_interfaces.TransactionHeaderRepository
	transactionDetailRepository   repository_interfaces.TransactionDetailRepository
	cafeProductRepository         repository_interfaces.CafeProductRepository
	cafeFranchiseRepository       repository_interfaces.CafeFranchiseRepository
	productRepository             repository_interfaces.ProductRepository
	cartRepository                repository_interfaces.CartRepository
	transaction                   repository_interfaces.Transaction
	productCache                  cache_interfaces.ProductCache
	midtransTransactionRepository repository_interfaces.MidtransTransactionRepository
}

func (t *TransactionConsumer) MidtransTransactionRecord(ctx context.Context) {
	for {
		//1. get transactionCode
		transactionCode, err := t.consumer.ReadMessage(ctx, enums.TransactionTopic, `transaction-code`)
		if err != nil {
			log.Println(`failed to get transactionCode`)
		}

		//2. get items checkouts
		itemCheckouts, err := t.consumer.ReadMessage(ctx, enums.TransactionTopic, fmt.Sprintf(`transaction-detail-%s`, transactionCode))
		if err != nil {
			log.Println(`error while try to consume transaction-detail`, err)
		}
		checkoutLists := utils.ReadJSONObject[[]*requests.CheckoutList](itemCheckouts)
		items := make([]midtrans.ItemDetails, len(*checkoutLists))

		var totalAmount int64 = 0
		transactionHeader, err := t.transactionRepository.FindByTransactionCode(ctx, transactionCode)
		if err != nil {
			log.Println(`failed to get transaction header`, err)
			return
		}

		cafeFranchise, err := t.consumer.ReadMessage(ctx, enums.TransactionTopic, fmt.Sprintf(`transaction-cafe-franchise-name-%s`, transactionHeader.TransactionCode))
		if err != nil {
			log.Println(apperror.ErrInternalServer, `failed to consume cafeFranchiseName`, err)
		}

		userIdStr, err := t.consumer.ReadMessage(ctx, enums.TransactionTopic, fmt.Sprintf(`transaction-cafe-user-%s`, transactionHeader.TransactionCode))
		if err != nil {
			log.Println(apperror.ErrInternalServer, `failed to consume userId`, err)
		}

		cafeIdStr, err := t.consumer.ReadMessage(ctx, enums.TransactionTopic, fmt.Sprintf(`transaction-cafe-cafe-%s`, transactionHeader.TransactionCode))
		if err != nil {
			log.Println(apperror.ErrInternalServer, `failed to consume cafeId`, err)
		}

		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			log.Println(`failed convert userId `, err)
			return
		}
		cafeId, err := strconv.Atoi(cafeIdStr)
		if err != nil {
			log.Println(`failed convert cafeId `, err)
			return
		}

		user, err := t.userRepository.FindByID(ctx, uint64(userId))
		if err != nil {
			log.Println(apperror.ErrInternalServer, `failed to get user id`, err)
		}

		customErr := t.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
			for index, item := range *checkoutLists {
				cafeProduct, err := t.cafeProductRepository.FindByID(ctx, item.CafeProductID)
				if err != nil {
					return apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe product`, err)
				}

				product, err := t.productRepository.FindByID(ctx, cafeProduct.ProductID)
				if err != nil {
					return apperror.NewCustomError(apperror.ErrInternalServer, `failed to get product`, err)
				}

				if cafeProduct.Stock < item.Qty {
					return apperror.NewCustomError(apperror.ErrBadRequest, `insufficient stock`, fmt.Errorf(`insufficient stock`))
				}

				cafeProduct.Stock -= item.Qty
				cafeProduct.UpdatedAt = time.Now()
				cafeProduct, err = t.cafeProductRepository.Update(ctx, cafeProduct)
				if err != nil {
					return apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe product`, err)
				}

				transactionDetail := &models.TransactionDetail{
					TransactionID: transactionHeader.ID,
					CafeProductID: item.CafeProductID,
					Qty:           item.Qty,
				}

				transactionDetail, err = t.transactionDetailRepository.Save(ctx, transactionDetail)
				if err != nil {
					return apperror.NewCustomError(apperror.ErrInternalServer, `failed to create detail transaction`, err)
				}

				items[index] = midtrans.ItemDetails{
					ID:           strconv.FormatUint(item.CafeProductID, 10),
					Name:         product.Name,
					Qty:          int32(item.Qty),
					Price:        int64(cafeProduct.Price),
					MerchantName: cafeFranchise,
				}

				totalAmount += int64(cafeProduct.Price * item.Qty)

				err = t.cartRepository.DeleteByUserIDAndCafeProductID(ctx, uint64(userId), item.CafeProductID)
				if err != nil {
					return apperror.NewCustomError(apperror.ErrInternalServer, `failed to delete cart`, err)
				}
				ctx = context.WithValue(ctx, enums.QtyTrx, item.Qty)
				_ = t.productCache.DecreaseProductStock(ctx, fmt.Sprintf(enums.CafeProductsKey, []uint64{uint64(cafeId)}), item.CafeProductID, enums.QtyTrx)
			}
			return nil
		})
		if customErr != nil {
			log.Println(`error trx`, customErr)
			return
		}
		names := strings.Split(user.Name, " ")
		midtransRequest := &midtrans_client.MidtransRequest{
			FName:           names[0],
			LName:           names[1],
			Email:           user.Email,
			Phone:           user.PhoneNumber,
			GrossAmt:        totalAmount,
			TransactionCode: transactionHeader.TransactionCode,
			Items:           items,
		}

		midtransResponse, err := t.midtrans.CreateTransaction(ctx, midtransRequest)
		if err != nil {
			log.Println(`error midtrans trans`, err)
			return
		}

		transactionHeader.Status = enums.TransactionSUCCESS.String()
		transactionHeader.UpdatedAt = time.Now()
		_, err = t.transactionRepository.Update(ctx, transactionHeader)
		if err != nil {
			log.Println(`failed to update transaction repository`, err)
			return
		}

		_, err = t.midtransTransactionRepository.Save(ctx, &models.MidtransTransaction{
			TransactionCode: midtransRequest.TransactionCode,
			TotalAmt:        midtransRequest.GrossAmt,
			Token:           midtransResponse.Token,
			RedirectURL:     midtransResponse.RedirectURL,
		})
		if err != nil {
			log.Println(`failed to create midtrans transaction`, err)
			return
		}

		err = t.publisher.WriteMessage(ctx, enums.TransactionTopic, fmt.Sprintf(`transaction-midtrans-response-%s`, transactionHeader.TransactionCode), utils.WriteJSONObject[midtrans_client.MidtransResponse](midtransResponse))
		if err != nil {
			log.Println(`error while try to publish transaction-midtrans-response`, err)
			return
		}
	}
}

var _ consumer_interfaces.TransactionConsumer = &TransactionConsumer{}

func newTransactionConsumer(consumer kafka_client.Consumer,
	publisher kafka_client.Publisher,
	midtrans midtrans_client.Midtrans,
	transactionRepository repository_interfaces.TransactionHeaderRepository,
	userRepository repository_interfaces.UserRepository,
	transactionDetailRepository repository_interfaces.TransactionDetailRepository,
	cafeProductRepository repository_interfaces.CafeProductRepository,
	cafeFranchiseRepository repository_interfaces.CafeFranchiseRepository,
	productRepository repository_interfaces.ProductRepository,
	cartRepository repository_interfaces.CartRepository,
	transaction repository_interfaces.Transaction,
	productCache cache_interfaces.ProductCache,
	midtransTransactionRepository repository_interfaces.MidtransTransactionRepository) *TransactionConsumer {
	return &TransactionConsumer{
		consumer:                      consumer,
		publisher:                     publisher,
		midtrans:                      midtrans,
		transactionRepository:         transactionRepository,
		userRepository:                userRepository,
		transactionDetailRepository:   transactionDetailRepository,
		cafeFranchiseRepository:       cafeFranchiseRepository,
		cartRepository:                cartRepository,
		cafeProductRepository:         cafeProductRepository,
		productRepository:             productRepository,
		transaction:                   transaction,
		productCache:                  productCache,
		midtransTransactionRepository: midtransTransactionRepository,
	}
}
