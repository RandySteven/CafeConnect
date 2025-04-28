package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type cafeProductRepository struct {
	dbx repository_interfaces.DBX
}

func (c *cafeProductRepository) Save(ctx context.Context, entity *models.CafeProduct) (result *models.CafeProduct, err error) {
	id, err := mysql_client.Save[models.CafeProduct](ctx, c.dbx(ctx), queries.InsertCafeProduct, &entity.CafeID, &entity.ProductID, &entity.Price)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (c *cafeProductRepository) FindByID(ctx context.Context, id uint64) (result *models.CafeProduct, err error) {
	result = &models.CafeProduct{}
	err = mysql_client.FindByID[models.CafeProduct](ctx, c.dbx(ctx), queries.SelectCafeProductByID, id, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *cafeProductRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.CafeProduct, err error) {
	return mysql_client.FindAll[models.CafeProduct](ctx, c.dbx(ctx), queries.SelectCafeProducts)
}

var _ repository_interfaces.CafeProductRepository = &cafeProductRepository{}

func newCafeProductRepository(dbx repository_interfaces.DBX) *cafeProductRepository {
	return &cafeProductRepository{
		dbx: dbx,
	}
}
