package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type cafeFranchiseRepository struct {
	dbx repository_interfaces.DBX
}

func (c *cafeFranchiseRepository) Save(ctx context.Context, entity *models.CafeFranchise) (result *models.CafeFranchise, err error) {
	id, err := mysql_client.Save[models.Cafe](ctx, c.dbx(ctx), queries.InsertIntoCafeFranchise, &entity.Name, &entity.LogoURL)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (c *cafeFranchiseRepository) FindByID(ctx context.Context, id uint64) (result *models.CafeFranchise, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *cafeFranchiseRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.CafeFranchise, err error) {
	return mysql_client.FindAll[models.CafeFranchise](ctx, c.dbx(ctx), queries.SelectCafeFranchises)
}

var _ repository_interfaces.CafeFranchiseRepository = &cafeFranchiseRepository{}

func newCafeFranchiseRepository(dbx repository_interfaces.DBX) *cafeFranchiseRepository {
	return &cafeFranchiseRepository{
		dbx: dbx,
	}
}
