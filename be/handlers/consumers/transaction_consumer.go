package consumers

import (
	"context"
	"errors"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apperror"
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
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

type TransactionConsumer struct {
	transactionTopic              topics_interfaces.TransactionTopic
	midtransTopic                 topics_interfaces.MidtransTopic
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

func (t *TransactionConsumer) MidtransTransactionRecord(ctx context.Context) error {
	return t.transactionTopic.RegisterConsumer(func(message string) {
		log.Println("Processing message:", message)

		transactionMessage := utils.ReadJSONObject[messages.TransactionMidtransMessage](message)
		transactionCode := transactionMessage.TransactionCode
		checkoutList := transactionMessage.CheckoutList

		items := make([]midtrans.ItemDetails, len(checkoutList))
		var totalAmount int64

		transactionHeader, err := t.transactionRepository.FindByTransactionCode(ctx, transactionCode)
		if err != nil {
			log.Println("failed to get transaction header:", err)
			return
		}

		for i, item := range checkoutList {
			cafeProduct, err := t.cafeProductRepository.FindByID(ctx, item.CafeProductID)
			if err != nil {
				log.Println("failed to find cafe product:", err)
				return
			}

			product, err := t.productRepository.FindByID(ctx, cafeProduct.ProductID)
			if err != nil {
				log.Println("failed to find product:", err)
				return
			}

			items[i] = midtrans.ItemDetails{
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

		midtransResponse, err := t.midtrans.CreateTransaction(ctx, midtransRequest)
		if err != nil {
			log.Println("error creating midtrans transaction:", err)
			return
		}

		_, err = t.midtransTransactionRepository.Save(ctx, &models.MidtransTransaction{
			TransactionCode: midtransRequest.TransactionCode,
			TotalAmt:        midtransRequest.GrossAmt,
			Token:           midtransResponse.Token,
			RedirectURL:     midtransResponse.RedirectURL,
		})
		if err != nil {
			log.Println("failed to save midtrans transaction:", err)
			return
		}

		_ = t.checkoutCache.SetMultiData(ctx, fmt.Sprintf(enums.TransactionCheckoutItemsKey, transactionCode), checkoutList)

		err = t.midtransTopic.WriteMessage(ctx, utils.WriteJSONObject[messages.MidtransPaymentConfirmationMessage](
			&messages.MidtransPaymentConfirmationMessage{
				Token:           midtransResponse.Token,
				TransactionCode: transactionHeader.TransactionCode,
			}))
		if err != nil {
			log.Println(`failed to publish to midtrans topic`, err)
			return
		}
	})
}

func (t *TransactionConsumer) MidtransPaymentConfirmation(ctx context.Context) error {
	return t.midtransTopic.RegisterConsumer(func(message string) {
		log.Println("Processing message:", message)

		midtransPaymentConfirm := utils.ReadJSONObject[messages.MidtransPaymentConfirmationMessage](message)

		err := utils.Retry(ctx, 10, func(ctx context.Context) error {
			return t.handleMidtransConfirmation(ctx, midtransPaymentConfirm)
		})

		if err != nil {
			log.Println("⚠️ Max retries reached or unrecoverable error:", err)
		}
	})
}

func (t *TransactionConsumer) handleMidtransConfirmation(ctx context.Context, confirm *messages.MidtransPaymentConfirmationMessage) error {
	log.Printf("Checking Midtrans status for %s", confirm.TransactionCode)

	transactionStatus, err := t.midtrans.CheckTransaction(ctx, confirm.TransactionCode)
	if err != nil {
		return fmt.Errorf("check transaction error: %w", err)
	}

	transactionHeader, err := t.transactionRepository.FindByTransactionCode(ctx, confirm.TransactionCode)
	if err != nil {
		return fmt.Errorf("get transaction header error: %w", err)
	}

	switch transactionStatus.TransactionStatus {
	case "settlement":
		log.Println("✅ Payment settled")
		transactionHeader.Status = enums.TransactionSUCCESS.String()
		transactionHeader.UpdatedAt = time.Now()

		checkoutList, err := t.checkoutCache.GetMultiData(ctx, fmt.Sprintf(enums.TransactionCheckoutItemsKey, transactionHeader.TransactionCode))
		if err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("redis error: %w", err)
		}

		customErr := t.transaction.RunInTx(ctx, func(ctx context.Context) *apperror.CustomError {
			for _, item := range checkoutList {
				cafeProduct, err := t.cafeProductRepository.FindByID(ctx, item.CafeProductID)
				if err != nil {
					return apperror.NewCustomError(apperror.ErrInternalServer, "failed to get cafe product", err)
				}

				if cafeProduct.Stock < item.Qty {
					return apperror.NewCustomError(apperror.ErrBadRequest, "insufficient stock", nil)
				}

				cafeProduct.Stock -= item.Qty
				cafeProduct.UpdatedAt = time.Now()

				_, err = t.cafeProductRepository.Update(ctx, cafeProduct)
				if err != nil {
					return apperror.NewCustomError(apperror.ErrInternalServer, "failed to update product", err)
				}

				transactionDetail := &models.TransactionDetail{
					TransactionID: transactionHeader.ID,
					CafeProductID: item.CafeProductID,
					Qty:           item.Qty,
				}

				_, err = t.transactionDetailRepository.Save(ctx, transactionDetail)
				if err != nil {
					return apperror.NewCustomError(apperror.ErrInternalServer, "failed to save transaction detail", err)
				}

				err = t.cartRepository.DeleteByUserIDAndCafeProductID(ctx, transactionHeader.UserID, item.CafeProductID)
				if err != nil {
					return apperror.NewCustomError(apperror.ErrInternalServer, "failed to delete cart", err)
				}
			}

			if _, err := t.transactionRepository.Update(ctx, transactionHeader); err != nil {
				return apperror.NewCustomError(apperror.ErrInternalServer, "failed to update transaction", err)
			}

			return nil
		})

		if customErr != nil {
			return fmt.Errorf("failed to commit transaction: %w", customErr)
		}

	case "pending":
		log.Println("ℹ️ Payment still pending")
		return fmt.Errorf("status still pending")

	case "cancel", "cancelled", "expire", "failure":
		log.Println("❌ Payment failed or expired:", transactionStatus.TransactionStatus)
		transactionHeader.Status = enums.TransactionFAILED.String()
		transactionHeader.UpdatedAt = time.Now()
		if _, err := t.transactionRepository.Update(ctx, transactionHeader); err != nil {
			return fmt.Errorf("update to failed status error: %w", err)
		}
	default:
		log.Println("⚠️ Unhandled Midtrans status:", transactionStatus.TransactionStatus)
	}

	return nil
}

var _ consumer_interfaces.TransactionConsumer = &TransactionConsumer{}

func newTransactionConsumer(
	transactionTopic topics_interfaces.TransactionTopic,
	midtransTopic topics_interfaces.MidtransTopic,
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
		midtransTopic:                 midtransTopic,
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
