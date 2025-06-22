package caches

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type checkoutCache struct {
	redis *redis.Client
}

func (c *checkoutCache) SetMultiData(ctx context.Context, key string, values []*requests.CheckoutList) (err error) {
	return redis_client.SetMultiple[requests.CheckoutList](ctx, c.redis, key, values)
}

func (c *checkoutCache) GetMultiData(ctx context.Context, key string) (values []*requests.CheckoutList, err error) {
	return redis_client.GetMultiple[requests.CheckoutList](ctx, c.redis, key)
}

var _ cache_interfaces.CheckoutCache = &checkoutCache{}

func newCheckoutCache(redis *redis.Client) *checkoutCache {
	return &checkoutCache{
		redis: redis,
	}
}
