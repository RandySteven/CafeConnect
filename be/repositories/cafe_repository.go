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
		entity.AddressID, entity.CafeFranchiseID, entity.CafeType, entity.PhotoURLs, entity.OpenHour, entity.CloseHour)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (c *cafeRepository) FindByID(ctx context.Context, id uint64) (result *models.Cafe, err error) {
	result = &models.Cafe{}
	err = mysql_client.FindByID[models.Cafe](ctx, c.dbx(ctx), queries.SelectCafeByID, id, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *cafeRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.Cafe, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *cafeRepository) FindByCafeFranchiseId(ctx context.Context, cafeFranchiseId uint64) (result []*models.Cafe, err error) {
	rows, err := c.dbx(ctx).QueryContext(ctx, queries.SelectCafesByCafeFranchiseID.String(), cafeFranchiseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cafe := &models.Cafe{}
		err = rows.Scan(
			&cafe.ID,
			&cafe.AddressID,
			&cafe.CafeFranchiseID,
			&cafe.CafeType,
			&cafe.PhotoURLs,
			&cafe.OpenHour,
			&cafe.CloseHour,
			&cafe.CreatedAt,
			&cafe.UpdatedAt,
			&cafe.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, cafe)
	}

	return result, nil
}

func (c *cafeRepository) FindByAddressId(ctx context.Context, addressId uint64) (result *models.Cafe, err error) {
	return
}

var _ repository_interfaces.CafeRepository = &cafeRepository{}

func newCafeRepository(dbx repository_interfaces.DBX) *cafeRepository {
	return &cafeRepository{
		dbx: dbx,
	}
}
