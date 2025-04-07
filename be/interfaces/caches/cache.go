package cache_interfaces

import "context"

type Cache[T any] interface {
	Set(ctx context.Context, request *T) (err error)
}
