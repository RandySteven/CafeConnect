package repositories

import (
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	"context"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type roleRepository struct {
	dbx repository_interfaces.DBX
}

func (r *roleRepository) Save(ctx context.Context, entity *models.Role) (result *models.Role, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) FindByID(ctx context.Context, id uint64) (result *models.Role, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.Role, err error) {
	return mysql_client.FindAll[models.Role](ctx, r.dbx(ctx), queries.SelectRolesQuery)
}

var _ repository_interfaces.RoleRepository = &roleRepository{}

func newRoleRepository(
	dbx repository_interfaces.DBX,
) *roleRepository {
	return &roleRepository{
		dbx: dbx,
	}
}
