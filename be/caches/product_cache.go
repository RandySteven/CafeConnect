package caches

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type productCache struct {
	redis *redis.Client
}

func (p *productCache) SetMultiData(ctx context.Context, key string, values []*responses.ListProductResponse) (err error) {
	return redis_client.SetMultiple[responses.ListProductResponse](ctx, p.redis, key, values)
}

func (p *productCache) GetMultiData(ctx context.Context, key string) (values []*responses.ListProductResponse, err error) {
	return redis_client.GetMultiple[responses.ListProductResponse](ctx, p.redis, key)
}

func (p *productCache) Set(ctx context.Context, key string, value *responses.DetailProductResponse) (err error) {
	return redis_client.Set[responses.DetailProductResponse](ctx, p.redis, key, value)
}

func (p *productCache) Get(ctx context.Context, key string) (value *responses.DetailProductResponse, err error) {
	return redis_client.Get[responses.DetailProductResponse](ctx, p.redis, key)
}

var _ cache_interfaces.ProductCache = &productCache{}

func newProductCache(redis *redis.Client) *productCache {
	return &productCache{
		redis: redis,
	}
}
