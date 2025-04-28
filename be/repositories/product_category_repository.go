package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type productCategoryRepository struct {
	dbx repository_interfaces.DBX
}

func (p *productCategoryRepository) FindByID(ctx context.Context, id uint64) (result *models.ProductCategory, err error) {
	result = &models.ProductCategory{}
	err = mysql_client.FindByID[models.ProductCategory](ctx, p.dbx(ctx), queries.SelectProductCategoriesByID, id, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *productCategoryRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.ProductCategory, err error) {
	return mysql_client.FindAll[models.ProductCategory](ctx, p.dbx(ctx), queries.SelectProductCategories)
}

var _ repository_interfaces.ProductCategoryRepository = &productCategoryRepository{}

func newProductCategoryRepository(dbx repository_interfaces.DBX) *productCategoryRepository {
	return &productCategoryRepository{
		dbx: dbx,
	}
}
