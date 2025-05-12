package caches

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type cafeCache struct {
	redis *redis.Client
}

func (c *cafeCache) SetCafeRadiusListCache(ctx context.Context, key string, response []*responses.ListCafeResponse) (err error) {
	return redis_client.SetMultiple[responses.ListCafeResponse](ctx, c.redis, key, response)
}

func (c *cafeCache) GetCafeRadiusListCache(ctx context.Context, key string) (response []*responses.ListCafeResponse, err error) {
	return redis_client.GetMultiple[responses.ListCafeResponse](ctx, c.redis, key)
}

func (c *cafeCache) SetCafeDetail(ctx context.Context, key string, value *responses.DetailCafeResponse) (err error) {
	return redis_client.Set[responses.DetailCafeResponse](ctx, c.redis, fmt.Sprintf(enums.CafeCacheKey, key), value)
}

func (c *cafeCache) GetCafeDetail(ctx context.Context, key string) (value *responses.DetailCafeResponse, err error) {
	return redis_client.Get[responses.DetailCafeResponse](ctx, c.redis, fmt.Sprintf(enums.CafeCacheKey, key))
}

func (c *cafeCache) GetFranchiseListCache(ctx context.Context) (result []*responses.FranchiseListResponse, err error) {
	return redis_client.GetMultiple[responses.FranchiseListResponse](ctx, c.redis, enums.FranchisesListCacheKey)
}

func (c *cafeCache) SetFranchiseListCache(ctx context.Context, response []*responses.FranchiseListResponse) (err error) {
	return redis_client.SetMultiple[responses.FranchiseListResponse](ctx, c.redis, enums.FranchisesListCacheKey, response)
}

var _ cache_interfaces.CafeCache = &cafeCache{}

func newCafeCache(redis *redis.Client) *cafeCache {
	return &cafeCache{
		redis: redis,
	}
}
