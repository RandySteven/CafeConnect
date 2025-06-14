package searcher_interfaces

import "context"

type (
	Indexer[T any] interface {
		SetIndex(ctx context.Context, index string, storeValue []*T)
		GetIndex(ctx context.Context, index string) (result []*T, err error)
	}

	Searcher[T any] interface {
		Search(ctx context.Context, index string, key string, value string) (result []*T, err error)
	}
)
