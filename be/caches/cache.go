package caches

import (
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	"github.com/redis/go-redis/v9"
)

type Caches struct {
	OnboardCache cache_interfaces.OnboardingCache
	CafeCache    cache_interfaces.CafeCache
}

func NewCaches(redis *redis.Client) *Caches {
	return &Caches{
		OnboardCache: newOnboardingCache(redis),
		CafeCache:    newCafeCache(redis),
	}
}
