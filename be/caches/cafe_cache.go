package caches

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type cafeCache struct {
	redis *redis.Client
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
