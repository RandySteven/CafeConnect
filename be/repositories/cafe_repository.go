package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type cafeRepository struct {
	dbx repository_interfaces.DBX
}

func (c *cafeRepository) Save(ctx context.Context, entity *models.Cafe) (result *models.Cafe, err error) {
	id, err := mysql_client.Save[models.Cafe](ctx, c.dbx(ctx), queries.InsertIntoCafe,
		entity.AddressID, entity.CafeFranchiseID, entity.CafeType, entity.PhotoURLs, entity.OpenHour.Format("15:04:04"), entity.CloseHour.Format("15:04:04"))
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (c *cafeRepository) FindByID(ctx context.Context, id uint64) (result *models.Cafe, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *cafeRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.Cafe, err error) {
	//TODO implement me
	panic("implement me")
}

var _ repository_interfaces.CafeRepository = &cafeRepository{}

func newCafeRepository(dbx repository_interfaces.DBX) *cafeRepository {
	return &cafeRepository{
		dbx: dbx,
	}
}
