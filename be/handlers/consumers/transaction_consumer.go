package consumers

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
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
		result, err := t.consumer.ReadMessage(ctx, `transaction`)
		if err != nil {
			log.Println(`failed to consumer result`, err)
			return
		}

		transactionMessage := utils.ReadJSONObject[messages.TransactionMidtransMessage](result)
		log.Println(`transaction message : `, transactionMessage)
		transactionCode := transactionMessage.TransactionCode
		checkoutList := transactionMessage.CheckoutList

		items := make([]midtrans.ItemDetails, len(checkoutList))
		var totalAmount int64 = 0
		transactionHeader, err := t.transactionRepository.FindByTransactionCode(ctx, transactionCode)
		if err != nil {
			log.Println(`failed to get transaction header`, err)
			return
		}
		log.Println(`success get transaction header `, transactionHeader)

		customErr := t.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic in transaction block: %v", r)
					customErr = apperror.NewCustomError(apperror.ErrInternalServer, "panic occurred", fmt.Errorf("%v", r))
				}
			}()
			for index, item := range checkoutList {
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
					MerchantName: transactionMessage.CafeFranchiseName,
				}

				totalAmount += int64(cafeProduct.Price * item.Qty)

				err = t.cartRepository.DeleteByUserIDAndCafeProductID(ctx, transactionMessage.UserID, item.CafeProductID)
				if err != nil {
					return apperror.NewCustomError(apperror.ErrInternalServer, `failed to delete cart`, err)
				}
			}
			return customErr
		})
		if customErr != nil {
			log.Println(`error trx`, customErr.Error())
			return
		}

		midtransRequest := &midtrans_client.MidtransRequest{
			FName:           transactionMessage.FName,
			LName:           transactionMessage.LName,
			Email:           transactionMessage.Email,
			Phone:           transactionMessage.Phone,
			GrossAmt:        totalAmount,
			TransactionCode: transactionHeader.TransactionCode,
			Items:           items,
		}
		log.Println(`midtrans request : `, midtransRequest)
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

		err = t.publisher.WriteMessage(ctx, fmt.Sprintf(`transaction-midtrans-response-%s`, transactionHeader.TransactionCode), utils.WriteJSONObject[midtrans_client.MidtransResponse](midtransResponse))
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
