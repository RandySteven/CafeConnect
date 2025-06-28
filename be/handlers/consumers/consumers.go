package consumers

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/caches"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	email_client "github.com/RandySteven/CafeConnect/be/pkg/email"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/RandySteven/CafeConnect/be/repositories"
	"github.com/RandySteven/CafeConnect/be/topics"
)

type (
	Consumers struct {
		//DummyConsumer       consumer_interfaces.DummyConsumer
		TransactionConsumer consumer_interfaces.TransactionConsumer
		OnboardingConsumer  consumer_interfaces.OnboardingConsumer
	}

	ConsumerFunc func(ctx context.Context) error

	Runners struct {
		ConsumerFunc []ConsumerFunc
	}
)

func consume(ctx context.Context, fn func(ctx context.Context)) {
	for {
		fn(ctx)
	}
}

func RegisterConsumer(consFunc ...ConsumerFunc) *Runners {
	return &Runners{
		ConsumerFunc: consFunc,
	}
}

func (r *Runners) Run(ctx context.Context) {
	for _, fun := range r.ConsumerFunc {
		go fun(ctx)
	}
}

func NewConsumers(
	repo *repositories.Repositories,
	cache *caches.Caches,
	topics *topics.Topics,
	midtrans midtrans_client.Midtrans,
	email email_client.Email,
) *Consumers {
	return &Consumers{
		//DummyConsumer: newDummyConsumer(topics),
		TransactionConsumer: newTransactionConsumer(
			topics.TransactionTopic,
			topics.MidtransTopic,
			topics.PointTopic,
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
			cache.CheckoutCache,
			repo.MidtransTransactionRepository),
		OnboardingConsumer: newOnboardingConsumer(
			topics.OnboardingTopic,
			topics.PointTopic,
			email,
			cache.OnboardCache,
			repo.UserRepository,
			repo.VerifyTokenRepository,
			repo.PointRepository),
	}
}
