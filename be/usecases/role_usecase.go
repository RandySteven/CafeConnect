package usecases

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
)

type roleUsecase struct {
	roleRepository repository_interfaces.RoleRepository
}

func (r *roleUsecase) GetRoleList(ctx context.Context) (result []*responses.RoleListResponse, customErr *apperror.CustomError) {
	roles, err := r.roleRepository.FindAll(ctx, 0, 0)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get roles`, err)
	}

	result = make([]*responses.RoleListResponse, len(roles))

	for index, role := range roles {
		result[index] = &responses.RoleListResponse{
			ID:        role.ID,
			Role:      role.Role,
			CreatedAt: role.CreatedAt,
			UpdatedAt: role.UpdatedAt,
			DeletedAt: role.DeletedAt,
		}
	}

	return result, nil
}

var _ usecase_interfaces.RoleUsecase = &roleUsecase{}

func newRoleUsecase(
	roleRepository repository_interfaces.RoleRepository,
) *roleUsecase {
	return &roleUsecase{
		roleRepository: roleRepository,
	}
}
