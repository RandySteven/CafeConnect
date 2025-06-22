package cache_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/payloads/requests"

type CheckoutCache interface {
	MultiDataCache[requests.CheckoutList]
}
