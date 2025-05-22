package cache_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/payloads/responses"

type ProductCache interface {
	MultiDataCache[responses.ListProductResponse]
	SingleDataCache[responses.DetailProductResponse]
}
