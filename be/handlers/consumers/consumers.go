package consumers

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/caches"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	email_client "github.com/RandySteven/CafeConnect/be/pkg/email"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/RandySteven/CafeConnect/be/repositories"
)

type Consumers struct {
	DummyConsumer       consumer_interfaces.DummyConsumer
	TransactionConsumer consumer_interfaces.TransactionConsumer
	OnboardingConsumer  consumer_interfaces.OnboardingConsumer
}

func consume(ctx context.Context, fn func(ctx context.Context)) {
	for {
		fn(ctx)
	}
}

func NewConsumers(
	repo *repositories.Repositories,
	cache *caches.Caches,
	consumer kafka_client.Consumer,
	publisher kafka_client.Publisher,
	midtrans midtrans_client.Midtrans,
	email email_client.Email,
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
		OnboardingConsumer: newOnboardingConsumer(
			consumer,
			publisher,
			email,
			repo.UserRepository),
	}
}
