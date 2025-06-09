package caches

import (
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	"github.com/redis/go-redis/v9"
)

type Caches struct {
	AddressCache     cache_interfaces.AddressCache
	OnboardCache     cache_interfaces.OnboardingCache
	CafeCache        cache_interfaces.CafeCache
	ProductCache     cache_interfaces.ProductCache
	TransactionCache cache_interfaces.TransactionCache
}

func NewCaches(redis *redis.Client) *Caches {
	return &Caches{
		AddressCache:     newAddressCache(redis),
		OnboardCache:     newOnboardingCache(redis),
		CafeCache:        newCafeCache(redis),
		ProductCache:     newProductCache(redis),
		TransactionCache: newTransactionCache(redis),
	}
}
