package cache_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type CafeCache interface {
	GetFranchiseListCache(ctx context.Context) (result []*responses.FranchiseListResponse, err error)
	SetFranchiseListCache(ctx context.Context, response []*responses.FranchiseListResponse) (err error)
	SetCafeDetail(ctx context.Context, key string, value *responses.DetailCafeResponse) (err error)
	GetCafeDetail(ctx context.Context, key string) (value *responses.DetailCafeResponse, err error)

	SetCafeRadiusListCache(ctx context.Context, key string, response []*responses.ListCafeResponse) (err error)
	GetCafeRadiusListCache(ctx context.Context, key string) (response []*responses.ListCafeResponse, err error)
}
