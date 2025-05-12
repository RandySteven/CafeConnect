package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type addressRepository struct {
	dbx repository_interfaces.DBX
}

func (a *addressRepository) Save(ctx context.Context, entity *models.Address) (result *models.Address, err error) {
	id, err := mysql_client.Save[models.Address](ctx, a.dbx(ctx), queries.InsertAddress, &entity.Address, &entity.Longitude, &entity.Latitude)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (a *addressRepository) FindByID(ctx context.Context, id uint64) (result *models.Address, err error) {
	result = &models.Address{}
	err = a.dbx(ctx).QueryRowContext(ctx, queries.SelectAddressByID.String(), id).Scan(
		&result.ID,
		&result.Address,
		&result.Longitude,
		&result.Latitude,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *addressRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.Address, err error) {
	return
}

func (a *addressRepository) Update(ctx context.Context, entity *models.Address) (result *models.Address, err error) {
	return
}

func (a *addressRepository) FindAddressBasedOnRadius(ctx context.Context, longitude, latitude float64, rangeKm uint64) (result []*models.Address, err error) {
	rows, err := a.dbx(ctx).QueryContext(ctx, queries.SelectAddressByRadiusNKm.String(), longitude, latitude, rangeKm)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		address := &models.Address{}
		err = rows.Scan(
			&address.ID,
			&address.Address,
			&address.Longitude,
			&address.Latitude,
			&address.CreatedAt,
			&address.UpdatedAt,
			&address.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, address)
	}

	return result, nil
}

var _ repository_interfaces.AddressRepository = &addressRepository{}

func newAddressRepository(dbx repository_interfaces.DBX) *addressRepository {
	return &addressRepository{
		dbx: dbx,
	}
}
