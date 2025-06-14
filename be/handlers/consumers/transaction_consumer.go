package consumers

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/enums"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/RandySteven/CafeConnect/be/utils"
	"log"
	"time"
)

type TransactionConsumer struct {
	consumer                      kafka_client.Consumer
	publisher                     kafka_client.Publisher
	midtrans                      midtrans_client.Midtrans
	transactionRepository         repository_interfaces.TransactionHeaderRepository
	midtransTransactionRepository repository_interfaces.MidtransTransactionRepository
}

func (t *TransactionConsumer) MidtransTransactionRecord(ctx context.Context) {
	//1. get midtrans request
	for {
		result, err := t.consumer.ReadMessage(ctx, enums.TransactionTopic, `transaction-midtrans-request`)
		if err != nil {
			log.Println(`error while try to consume transaction-midtrans-request`, err)
			return
		}
		log.Println(result)
		midtransRequest := utils.ReadJSONObject[midtrans_client.MidtransRequest](result)

		midtransResponse, err := t.midtrans.CreateTransaction(ctx, midtransRequest)
		if err != nil {
			log.Println(`error midtrans trans`, err)
			return
		}

		transactionHeader, _ := t.transactionRepository.FindByTransactionCode(ctx, midtransRequest.TransactionCode)
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
	midtransTransactionRepository repository_interfaces.MidtransTransactionRepository) *TransactionConsumer {
	return &TransactionConsumer{
		consumer:                      consumer,
		publisher:                     publisher,
		midtrans:                      midtrans,
		transactionRepository:         transactionRepository,
		midtransTransactionRepository: midtransTransactionRepository,
	}
}
