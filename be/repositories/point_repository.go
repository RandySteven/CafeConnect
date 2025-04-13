package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type pointRepository struct {
	dbx repository_interfaces.DBX
}

func (p *pointRepository) Save(ctx context.Context, entity *models.Point) (result *models.Point, err error) {
	id, err := mysql_client.Save[models.Point](ctx, p.dbx(ctx), queries.InsertPoint, &entity.Point, &entity.UserID)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (p *pointRepository) FindByID(ctx context.Context, id uint64) (result *models.Point, err error) {
	return
}

func (p *pointRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.Point, err error) {
	return
}

func (p *pointRepository) Update(ctx context.Context, entity *models.Point) (result *models.Point, err error) {
	return
}

func (p *pointRepository) FindByUserID(ctx context.Context, userId uint64) (result *models.Point, err error) {
	return
}

var _ repository_interfaces.PointRepository = &pointRepository{}

func newPointRepository(dbx repository_interfaces.DBX) *pointRepository {
	return &pointRepository{
		dbx: dbx,
	}
}
