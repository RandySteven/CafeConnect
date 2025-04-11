package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
)

type pointRepository struct {
	dbx repository_interfaces.DBX
}

func (p *pointRepository) Save(ctx context.Context, entity *models.Point) (result *models.Point, err error) {
	return nil, nil
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

var _ repository_interfaces.PointRepository = &pointRepository{}

func newPointRepository(dbx repository_interfaces.DBX) *pointRepository {
	return &pointRepository{
		dbx: dbx,
	}
}
