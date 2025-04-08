package redis_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/enums"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
)

var (
	redisTimeout = os.Getenv("REDIS_EXPIRATION")
	client       *redis.Client
	limiter      *redis_rate.Limiter
	rateLimiter  = os.Getenv("RATE_LIMITER")
)

type (
	Redis interface {
		Ping() error
		Client() *redis.Client
		ClearCache(ctx context.Context) error
	}

	redisClient struct {
		client  *redis.Client
		limiter *redis_rate.Limiter
	}
)

func NewRedisClient(config *configs.Config) (*redisClient, error) {
	redisCfg := config.Config.Redis
	addr := fmt.Sprintf("%s:%s", redisCfg.Host, redisCfg.Port)
	log.Println("connecting to redis : ", addr)

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisCfg.Password,
		DB:       0,
	})
	limiter = redis_rate.NewLimiter(client)
	return &redisClient{
		client:  client,
		limiter: limiter,
	}, nil
}

func (c *redisClient) Ping() error {
	return c.client.Ping(context.Background()).Err()
}

func (c *redisClient) Client() *redis.Client {
	return c.client
}

func (c *redisClient) ClearCache(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

func RateLimiter(ctx context.Context) error {
	rateLimiterInt, _ := strconv.Atoi(rateLimiter)
	clientIP := ctx.Value(enums.ClientIP).(string)
	res, err := limiter.Allow(ctx, clientIP, redis_rate.PerMinute(rateLimiterInt))
	if err != nil {
		return err
	}
	if res.Remaining == 0 {
		return errors.New("Rate limit exceeded")
	}
	return nil
}
