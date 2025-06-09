package caches

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type addressCache struct {
	redis *redis.Client
}

func (a *addressCache) SetMultiData(ctx context.Context, key string, values []*responses.GetUserAddressResponse) (err error) {
	return redis_client.SetMultiple[responses.GetUserAddressResponse](ctx, a.redis, key, values)
}

func (a *addressCache) GetMultiData(ctx context.Context, key string) (values []*responses.GetUserAddressResponse, err error) {
	return redis_client.GetMultiple[responses.GetUserAddressResponse](ctx, a.redis, key)
}

var _ cache_interfaces.AddressCache = &addressCache{}

func newAddressCache(redis *redis.Client) *addressCache {
	return &addressCache{
		redis: redis,
	}
}
