package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type cartRepository struct {
	dbx repository_interfaces.DBX
}

func (c *cartRepository) Update(ctx context.Context, entity *models.Cart) (result *models.Cart, err error) {
	err = mysql_client.Update[models.Cart](ctx, c.dbx(ctx), queries.UpdateCartByID,
		&entity.UserID, &entity.CafeProductID, &entity.Qty, &entity.CreatedAt, &entity.UpdatedAt, &entity.DeletedAt, &entity.ID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (c *cartRepository) Save(ctx context.Context, entity *models.Cart) (result *models.Cart, err error) {
	id, err := mysql_client.Save[models.Cart](ctx, c.dbx(ctx), queries.InsertIntoCart, &entity.UserID, &entity.CafeProductID, &entity.Qty)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (c *cartRepository) FindByUserID(ctx context.Context, userId uint64) (result []*models.Cart, err error) {
	rows, err := c.dbx(ctx).QueryContext(ctx, queries.SelectCartsByUserID.String(), userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		cart := &models.Cart{}
		err = rows.Scan(
			&cart.ID,
			&cart.UserID,
			&cart.CafeProductID,
			&cart.Qty,
			&cart.CreatedAt,
			&cart.UpdatedAt,
			&cart.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, cart)
	}

	return result, nil
}

func (c *cartRepository) FindByUserIDAndCafeProductID(ctx context.Context, userId uint64, cafeProductId uint64) (result *models.Cart, err error) {
	result = &models.Cart{}
	err = c.dbx(ctx).QueryRowContext(ctx, queries.SelectCartByUserIDAndCafeProductID.String(), userId, cafeProductId).Scan(
		&result.ID,
		&result.UserID,
		&result.CafeProductID,
		&result.Qty,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *cartRepository) DeleteByUserID(ctx context.Context, userId uint64) (err error) {
	_, err = c.dbx(ctx).ExecContext(ctx, queries.DeleteCartByUserID.String(), userId)
	if err != nil {
		return err
	}
	return nil
}

var _ repository_interfaces.CartRepository = &cartRepository{}

func newCartRepository(dbx repository_interfaces.DBX) *cartRepository {
	return &cartRepository{
		dbx: dbx,
	}
}
