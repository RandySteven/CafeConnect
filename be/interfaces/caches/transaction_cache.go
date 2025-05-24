package cache_interfaces

import "github.com/RandySteven/CafeConnect/be/entities/payloads/responses"

type TransactionCache interface {
	SingleDataCache[responses.TransactionDetailResponse]
	MultiDataCache[responses.TransactionListResponse]
}
