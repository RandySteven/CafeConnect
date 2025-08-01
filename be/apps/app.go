package apps

import (
	"context"
	caches2 "github.com/RandySteven/CafeConnect/be/caches"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/handlers/apis"
	consumers2 "github.com/RandySteven/CafeConnect/be/handlers/consumers"
	aws_client "github.com/RandySteven/CafeConnect/be/pkg/aws"
	cron_client "github.com/RandySteven/CafeConnect/be/pkg/cron"
	elastic_client "github.com/RandySteven/CafeConnect/be/pkg/elastic"
	email_client "github.com/RandySteven/CafeConnect/be/pkg/email"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	nsq_client "github.com/RandySteven/CafeConnect/be/pkg/nsq"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
	repositories2 "github.com/RandySteven/CafeConnect/be/repositories"
	topics2 "github.com/RandySteven/CafeConnect/be/topics"
	usecases2 "github.com/RandySteven/CafeConnect/be/usecases"
	"log"
)

type App struct {
	MySQL     mysql_client.MySQL
	Redis     redis_client.Redis
	Scheduler cron_client.Scheduler
	AWS       aws_client.AWS
	Midtrans  midtrans_client.Midtrans
	Kafka     kafka_client.Kafka
	Pub       kafka_client.Publisher
	Sub       kafka_client.Consumer
	Nsq       nsq_client.Nsq
	Elastic   elastic_client.Elastic
	Email     email_client.Email
}

func NewApps(config *configs.Config) (*App, error) {

	mysql, err := mysql_client.NewMySQLClient(config)
	if err != nil {
		return nil, err
	}

	redis, err := redis_client.NewRedisClient(config)
	if err != nil {
		return nil, err
	}

	scheduler, err := cron_client.NewScheduler(config)
	if err != nil {
		return nil, err
	}

	aws, err := aws_client.NewAWS(config)
	if err != nil {
		return nil, err
	}

	midtrans, err := midtrans_client.NewMidtrans(config)
	if err != nil {
		return nil, err
	}

	kafka, err := kafka_client.NewKafkaClient(config)
	if err != nil {
		kafka = nil
	}

	pub, err := kafka_client.NewPublisher(config)
	if err != nil {
		pub = nil
	}

	sub, err := kafka_client.NewConsumer(config)
	if err != nil {
		sub = nil
	}

	es, err := elastic_client.NewElastic(config)
	if err != nil {
		log.Println(`error es : `, err)
		es = nil
	}

	em, err := email_client.NewEmail(config)
	if err != nil {
		return nil, err
	}

	nsq, err := nsq_client.NewNsqClient(config)
	if err != nil {
		return nil, err
	}

	return &App{
		MySQL:     mysql,
		Redis:     redis,
		Scheduler: scheduler,
		AWS:       aws,
		Midtrans:  midtrans,
		Kafka:     kafka,
		Pub:       pub,
		Sub:       sub,
		Nsq:       nsq,
		Elastic:   es,
		Email:     em,
	}, nil
}

func (a *App) PrepareHttpHandler(ctx context.Context) *apis.APIs {
	repositories := repositories2.NewRepositories(a.MySQL.Client())
	caches := caches2.NewCaches(a.Redis.Client())
	topics := topics2.NewTopics(a.Nsq)
	usecases := usecases2.NewUsecases(repositories, caches, a.AWS, topics, a.Midtrans)
	return apis.NewAPIs(usecases)
}

func (a *App) RefreshRedis(ctx context.Context) error {
	return a.Redis.ClearCache(ctx)
}

func (a *App) PrepareJobScheduler(ctx context.Context) {
}

func (a *App) PrepareConsumer(ctx context.Context) *consumers2.Consumers {
	topics := topics2.NewTopics(a.Nsq)
	repositories := repositories2.NewRepositories(a.MySQL.Client())
	caches := caches2.NewCaches(a.Redis.Client())
	consumers := consumers2.NewConsumers(repositories, caches, topics, a.Midtrans, a.Email)
	return consumers
}
