package consumers

import (
	"github.com/RandySteven/CafeConnect/be/caches"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/RandySteven/CafeConnect/be/repositories"
)

type Consumers struct {
	DummyConsumer       consumer_interfaces.DummyConsumer
	TransactionConsumer consumer_interfaces.TransactionConsumer
}

func NewConsumers(
	repo *repositories.Repositories,
	cache *caches.Caches,
	consumer kafka_client.Consumer,
	publisher kafka_client.Publisher,
	midtrans midtrans_client.Midtrans,
) *Consumers {
	return &Consumers{
		DummyConsumer: newDummyConsumer(consumer),
		TransactionConsumer: newTransactionConsumer(
			consumer,
			publisher,
			midtrans,
			repo.TransactionHeaderRepository,
			repo.UserRepository,
			repo.TransactionDetailRepository,
			repo.CafeProductRepository,
			repo.CafeFranchiseRepository,
			repo.ProductRepository,
			repo.CartRepository,
			repo.Transaction,
			cache.ProductCache,
			repo.MidtransTransactionRepository),
	}
}
