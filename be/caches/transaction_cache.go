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

type transactionCache struct {
	redis *redis.Client
}

func (t *transactionCache) Set(ctx context.Context, key string, value *responses.TransactionDetailResponse) (err error) {
	return redis_client.Set[responses.TransactionDetailResponse](ctx, t.redis, fmt.Sprintf(enums.TransactionDetailKey, key), value)
}

func (t *transactionCache) Get(ctx context.Context, key string) (value *responses.TransactionDetailResponse, err error) {
	return redis_client.Get[responses.TransactionDetailResponse](ctx, t.redis, fmt.Sprintf(enums.TransactionDetailKey, key))
}

func (t *transactionCache) SetMultiData(ctx context.Context, key string, values []*responses.TransactionListResponse) (err error) {
	return redis_client.SetMultiple[responses.TransactionListResponse](ctx, t.redis, fmt.Sprintf(enums.TransactionListKey, key), values)
}

func (t *transactionCache) GetMultiData(ctx context.Context, key string) (values []*responses.TransactionListResponse, err error) {
	return redis_client.GetMultiple[responses.TransactionListResponse](ctx, t.redis, fmt.Sprintf(enums.TransactionListKey, key))
}

var _ cache_interfaces.TransactionCache = &transactionCache{}

func newTransactionCache(redis *redis.Client) *transactionCache {
	return &transactionCache{
		redis: redis,
	}
}
