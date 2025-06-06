package cache_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type ProductCache interface {
	MultiDataCache[responses.ListProductResponse]
	SingleDataCache[responses.DetailProductResponse]
	DecreaseProductStock(ctx context.Context, key string, productId uint64, triggerCtx string) (err error)
}
