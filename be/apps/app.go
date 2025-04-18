package apps

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/handlers/apis"
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

	return &App{
		MySQL:         mysql,
		Redis:         redis,
		GoogleStorage: googleStorage,
	}, nil
}

func (a *App) PrepareHttpHandler(ctx context.Context) *apis.APIs {
	repositories := repositories2.NewRepositories(a.MySQL.Client())
	usecases := usecases2.NewUsecases(repositories)
	return apis.NewAPIs(usecases)
}

func (a *App) PrepareJobScheduler(ctx context.Context) {

}
