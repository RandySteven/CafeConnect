package repositories

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
	"github.com/RandySteven/CafeConnect/be/utils"
	"log"
)

type cafeProductRepository struct {
	dbx repository_interfaces.DBX
}

func (c *cafeProductRepository) Update(ctx context.Context, entity *models.CafeProduct) (result *models.CafeProduct, err error) {
	err = mysql_client.Update[models.CafeProduct](ctx, c.dbx(ctx), queries.UpdateCafeProductByID,
		&entity.CafeID,
		&entity.ProductID,
		&entity.Price,
		&entity.Stock,
		&entity.Status,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.ID,
	)
	if err != nil {
		return nil, err
	}
	return entity, nil
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

func (c *cafeProductRepository) FindByCafeID(ctx context.Context, cafeID uint64) (result []*models.CafeProduct, err error) {
	rows, err := c.dbx(ctx).QueryContext(ctx, queries.SelectCafeProductsByCafeID.String(), cafeID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		res := &models.CafeProduct{}
		err = rows.Scan(&res.ID, &res.CafeID, &res.ProductID, &res.Price, &res.Stock, &res.Status,
			&res.CreatedAt, &res.UpdatedAt, &res.DeletedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}

	return result, nil
}

func (c *cafeProductRepository) FindByCafeIDs(ctx context.Context, cafeIDs []uint64) (result []*models.CafeProduct, err error) {
	query := queries.SelectCafeProductsInCafeIDs.String() + utils.InQuery(cafeIDs)
	log.Println(query)
	rows, err := c.dbx(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		res := &models.CafeProduct{}
		err = rows.Scan(&res.ID, &res.CafeID, &res.ProductID, &res.Price, &res.Stock, &res.Status,
			&res.CreatedAt, &res.UpdatedAt, &res.DeletedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}

func (c *cafeProductRepository) FindCafeIDByCafeProductIDs(ctx context.Context, cafeProductIDs []uint64) (cafeId uint64, err error) {
	query := fmt.Sprintf(queries.SelectCafeIdByCafeProductIDs.String(), utils.InQuery(cafeProductIDs))
	err = c.dbx(ctx).QueryRowContext(ctx, query).Scan(cafeId)
	if err != nil {
		return cafeId, err
	}
	return cafeId, nil
}

var _ repository_interfaces.CafeProductRepository = &cafeProductRepository{}

func newCafeProductRepository(dbx repository_interfaces.DBX) *cafeProductRepository {
	return &cafeProductRepository{
		dbx: dbx,
	}
}
