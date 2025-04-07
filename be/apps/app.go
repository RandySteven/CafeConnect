package apps

import (
	"github.com/RandySteven/CafeConnect/be/configs"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
)

type App struct {
	MySQL mysql_client.MySQL
	Redis redis_client.Redis
}

func NewApps(config *configs.Config) (*App, error) {
	return &App{}, nil
}
