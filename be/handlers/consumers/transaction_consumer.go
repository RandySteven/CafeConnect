package consumers

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/midtrans/midtrans-go"
	"log"
	"strconv"
)

type TransactionConsumer struct {
	transactionTopic              topics_interfaces.TransactionTopic
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
	checkoutCache                 cache_interfaces.CheckoutCache
	midtransTransactionRepository repository_interfaces.MidtransTransactionRepository
}

func (t *TransactionConsumer) MidtransTransactionRecord(ctx context.Context) {
	for {
		result, err := t.transactionTopic.ReadMessage(ctx, `transaction`)
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

		for index, item := range checkoutList {
			cafeProduct, err := t.cafeProductRepository.FindByID(ctx, item.CafeProductID)
			if err != nil {
				return
			}

			product, err := t.productRepository.FindByID(ctx, cafeProduct.ProductID)
			if err != nil {
				return
			}

			items[index] = midtrans.ItemDetails{
				ID:           strconv.FormatUint(item.CafeProductID, 10),
				Name:         product.Name,
				Qty:          int32(item.Qty),
				Price:        int64(cafeProduct.Price),
				MerchantName: transactionMessage.CafeFranchiseName,
			}

			totalAmount += int64(cafeProduct.Price * item.Qty)
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

		//save to midtrans_transactions table
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

		_ = t.checkoutCache.SetMultiData(ctx, fmt.Sprintf(enums.TransactionCheckoutItemsKey, transactionCode), checkoutList)

		err = t.transactionTopic.WriteMessage(ctx, fmt.Sprintf(`transaction-midtrans-response-%s`, transactionHeader.TransactionCode), utils.WriteJSONObject[midtrans_client.MidtransResponse](midtransResponse))
		if err != nil {
			log.Println(`error while try to publish transaction-midtrans-response`, err)
			return
		}
	}
}

var _ consumer_interfaces.TransactionConsumer = &TransactionConsumer{}

func newTransactionConsumer(
	transactionTopic topics_interfaces.TransactionTopic,
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
	checkoutCache cache_interfaces.CheckoutCache,
	midtransTransactionRepository repository_interfaces.MidtransTransactionRepository) *TransactionConsumer {
	return &TransactionConsumer{
		transactionTopic:              transactionTopic,
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
		checkoutCache:                 checkoutCache,
		midtransTransactionRepository: midtransTransactionRepository,
	}
}
