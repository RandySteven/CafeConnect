package caches

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
	"github.com/redis/go-redis/v9"
	"log"
)

type productCache struct {
	redis *redis.Client
}

func (p *productCache) Del(ctx context.Context, key string) (err error) {
	return redis_client.Del(ctx, p.redis, key)
}

func (p *productCache) SetMultiData(ctx context.Context, key string, values []*responses.ListProductResponse) (err error) {
	return redis_client.SetMultiple[responses.ListProductResponse](ctx, p.redis, key, values)
}

func (p *productCache) GetMultiData(ctx context.Context, key string) (values []*responses.ListProductResponse, err error) {
	return redis_client.GetMultiple[responses.ListProductResponse](ctx, p.redis, key)
}

func (p *productCache) Set(ctx context.Context, key string, value *responses.DetailProductResponse) (err error) {
	return redis_client.Set[responses.DetailProductResponse](ctx, p.redis, key, value)
}

func (p *productCache) Get(ctx context.Context, key string) (value *responses.DetailProductResponse, err error) {
	return redis_client.Get[responses.DetailProductResponse](ctx, p.redis, key)
}

func (p *productCache) DecreaseProductStock(ctx context.Context, key string, productId uint64, triggerCtx string) (err error) {
	log.Println(key)
	products, err := p.GetMultiData(ctx, key)
	if err != nil {
		return err
	}

	var qty uint64 = 0

	updated := false
	for _, product := range products {
		if product.ID == productId {
			switch triggerCtx {
			case enums.QtyCart:
				qty = ctx.Value(enums.QtyCart).(uint64)
				_ = redis_client.Set[uint64](ctx, p.redis, `qty_cart`, &qty)

				product.Stock -= qty
			case enums.QtyTrx:
				abs := 1
				cartQty, _ := redis_client.Get[uint64](ctx, p.redis, `qty_cart`)
				temp := int(ctx.Value(enums.QtyTrx).(uint64)) - int(*cartQty)
				if ctx.Value(enums.QtyTrx).(uint64) < *cartQty {
					abs *= -1
				}
				temp *= abs
				qty = uint64(temp)
				if abs == -1 {
					product.Stock += qty
				} else {
					product.Stock -= qty
				}
				_ = redis_client.Del[uint64](ctx, p.redis, `qty_cart`)
			}
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf(`failed to update product stock`)
	}

	return p.SetMultiData(ctx, key, products)
}

var _ cache_interfaces.ProductCache = &productCache{}

func newProductCache(redis *redis.Client) *productCache {
	return &productCache{
		redis: redis,
	}
}
