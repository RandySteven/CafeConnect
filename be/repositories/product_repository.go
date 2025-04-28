package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type productRepository struct {
	dbx repository_interfaces.DBX
}

func (p *productRepository) Save(ctx context.Context, entity *models.Product) (result *models.Product, err error) {
	id, err := mysql_client.Save[models.Product](ctx, p.dbx(ctx), queries.InsertIntoProduct, &entity.Name, &entity.PhotoURL, &entity.ProductCategoryID)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (p *productRepository) FindByID(ctx context.Context, id uint64) (result *models.Product, err error) {
	result = &models.Product{}
	err = mysql_client.FindByID[models.Product](ctx, p.dbx(ctx), queries.SelectProductByID, id, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *productRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.Product, err error) {
	return mysql_client.FindAll[models.Product](ctx, p.dbx(ctx), queries.SelectProducts)
}

var _ repository_interfaces.ProductRepository = &productRepository{}

func newProductRepository(dbx repository_interfaces.DBX) *productRepository {
	return &productRepository{
		dbx: dbx,
	}
}
