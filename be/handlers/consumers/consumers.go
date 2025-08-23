package consumers

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/caches"
	consumer_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/consumers"
	email_client "github.com/RandySteven/CafeConnect/be/pkg/email"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	nsq_client "github.com/RandySteven/CafeConnect/be/pkg/nsq"
	"github.com/RandySteven/CafeConnect/be/repositories"
	"github.com/RandySteven/CafeConnect/be/topics"
	"log"
)

type (
	Consumers struct {
		//DummyConsumer       consumer_interfaces.DummyConsumer
		TransactionConsumer consumer_interfaces.TransactionConsumer
		OnboardingConsumer  consumer_interfaces.OnboardingConsumer
	}

	ConsumerFunc func(ctx context.Context) error

	RunConsumer map[string]ConsumerFunc

	Runners struct {
		nsq          nsq_client.Nsq
		ConsumerFunc []ConsumerFunc
		RunConsumers RunConsumer
	}
)

func consume(ctx context.Context, fn func(ctx context.Context)) {
	for {
		fn(ctx)
	}
}

func InitRunner(nsq nsq_client.Nsq) *Runners {
	return &Runners{
		nsq:          nsq,
		RunConsumers: make(map[string]ConsumerFunc),
	}
}

func (r *Runners) RegisterConsumer(topic string, fun ConsumerFunc) *Runners {
	r.RunConsumers[topic] = fun
	return r
}

func (r *Runners) Run(ctx context.Context) error {
	errChan := make(chan error, len(r.RunConsumers))

	for topic, consumer := range r.RunConsumers {
		go func(topic string, consumer ConsumerFunc) {
			log.Println(`execute consumer `, consumer)
			err := r.nsq.RegisterConsumer(topic, func(msgCtx context.Context, key string) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic in consumer %s: %v", topic, r)
					}
				}()

				if err := consumer(msgCtx); err != nil {
					log.Printf("Error in consumer %s: %v", topic, err)
				}
			})
			if err != nil {
				errChan <- fmt.Errorf("failed to register consumer for topic %s: %w", topic, err)
			}
		}(topic, consumer)
	}

	select {
	case err := <-errChan:
		return err
	default:
		return nil
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
