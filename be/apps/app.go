package apps

import (
	"context"
	caches2 "github.com/RandySteven/CafeConnect/be/caches"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/handlers/apis"
	aws_client "github.com/RandySteven/CafeConnect/be/pkg/aws"
	cron_client "github.com/RandySteven/CafeConnect/be/pkg/cron"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
	storage_client "github.com/RandySteven/CafeConnect/be/pkg/storage"
	repositories2 "github.com/RandySteven/CafeConnect/be/repositories"
	usecases2 "github.com/RandySteven/CafeConnect/be/usecases"
)

type App struct {
	MySQL         mysql_client.MySQL
	Redis         redis_client.Redis
	GoogleStorage storage_client.GoogleStorage
	Scheduler     cron_client.Scheduler
	AWS           aws_client.AWS
	Midtrans      midtrans_client.Midtrans
	Kafka         kafka_client.Kafka
	Pub           kafka_client.Publisher
	Sub           kafka_client.Consumer
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

	googleStorage, err := storage_client.NewGoogleStorage(config)
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

	return &App{
		MySQL:         mysql,
		Redis:         redis,
		GoogleStorage: googleStorage,
		Scheduler:     scheduler,
		AWS:           aws,
		Midtrans:      midtrans,
		Kafka:         kafka,
		Pub:           pub,
		Sub:           sub,
	}, nil
}

func (a *App) PrepareHttpHandler(ctx context.Context) *apis.APIs {
	repositories := repositories2.NewRepositories(a.MySQL.Client())
	caches := caches2.NewCaches(a.Redis.Client())
	usecases := usecases2.NewUsecases(repositories, caches, a.GoogleStorage, a.AWS, a.Pub, a.Sub, a.Midtrans)
	return apis.NewAPIs(usecases)
}

func (a *App) RefreshRedis(ctx context.Context) error {
	return a.Redis.ClearCache(ctx)
}

func (a *App) PrepareJobScheduler(ctx context.Context) {
}

func (a *App) PrepareConsumer(ctx context.Context) {
}
