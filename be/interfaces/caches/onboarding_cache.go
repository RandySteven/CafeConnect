package cache_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type OnboardingCache interface {
	SingleDataCache[responses.OnboardUserResponse]
	Del(ctx context.Context, key string) error
}
