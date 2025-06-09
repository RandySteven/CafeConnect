package usecase_interfaces

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
)

type RoleUsecase interface {
	GetRoleList(ctx context.Context) (result []*responses.RoleListResponse, customErr *apperror.CustomError)
}
