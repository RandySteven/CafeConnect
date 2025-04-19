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

type onboardingCache struct {
	redis *redis.Client
}

func (o *onboardingCache) Set(ctx context.Context, key string, value *responses.OnboardUserResponse) (err error) {
	return redis_client.Set[responses.OnboardUserResponse](ctx, o.redis, fmt.Sprintf(enums.OnboardUserCacheKey, key), value)
}

func (o *onboardingCache) Get(ctx context.Context, key string) (value *responses.OnboardUserResponse, err error) {
	return redis_client.Get[responses.OnboardUserResponse](ctx, o.redis, fmt.Sprintf(enums.OnboardUserCacheKey, key))
}

var _ cache_interfaces.OnboardingCache = &onboardingCache{}

func newOnboardingCache(redis *redis.Client) *onboardingCache {
	return &onboardingCache{
		redis: redis,
	}
}
