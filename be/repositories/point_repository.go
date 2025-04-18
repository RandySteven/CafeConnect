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
	result = &models.Point{}
	err = mysql_client.FindByID[models.Point](ctx, p.dbx(ctx), queries.SelectPointByID, id, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *pointRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.Point, err error) {
	return mysql_client.FindAll[models.Point](ctx, p.dbx(ctx), queries.SelectPoints)
}

func (p *pointRepository) Update(ctx context.Context, entity *models.Point) (result *models.Point, err error) {
	err = mysql_client.Update[models.Point](ctx, p.dbx(ctx), queries.UpdatePointByID, &entity.Point, &entity.UserID, &entity.CreatedAt, &entity.UpdatedAt, &entity.DeletedAt, &entity.ID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (p *pointRepository) FindByUserID(ctx context.Context, userId uint64) (result *models.Point, err error) {
	result = &models.Point{}
	err = p.dbx(ctx).QueryRowContext(ctx, queries.SelectPointByUserID.String(), userId).Scan(
		&result.ID,
		&result.Point,
		&result.UserID,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

var _ repository_interfaces.PointRepository = &pointRepository{}

func newPointRepository(dbx repository_interfaces.DBX) *pointRepository {
	return &pointRepository{
		dbx: dbx,
	}
}
