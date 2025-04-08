package apps

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/configs"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
)

type App struct {
	MySQL mysql_client.MySQL
	Redis redis_client.Redis
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

	return &App{
		MySQL: mysql,
		Redis: redis,
	}, nil
}

func (a *App) PrepareHttpHandler(ctx context.Context) {

}

func (a *App) PrepareJobScheduler(ctx context.Context) {

}
