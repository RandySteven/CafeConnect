package cache_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/payloads/responses"

type OnboardingCache interface {
	SingleDataCache[responses.OnboardUserResponse]
}
