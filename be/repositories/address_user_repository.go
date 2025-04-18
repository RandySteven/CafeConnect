package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type addressUserRepository struct {
	dbx repository_interfaces.DBX
}

func (a *addressUserRepository) FindByAddressAndUserID(ctx context.Context, addressID uint64, userID uint64) (result *models.AddressUser, err error) {
	result = &models.AddressUser{}
	err = a.dbx(ctx).QueryRowContext(ctx, queries.SelectAddressUserByAddressAndUserID.String(), addressID, userID).
		Scan(
			&result.ID,
			&result.AddressID,
			&result.UserID,
			&result.IsDefault,
			&result.CreatedAt,
			&result.UpdatedAt,
			&result.DeletedAt,
		)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *addressUserRepository) FindByUserID(ctx context.Context, userID uint64) (results []*models.AddressUser, err error) {
	rows, err := a.dbx(ctx).QueryContext(ctx, queries.SelectAddressUserByUserID.String(), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		result := &models.AddressUser{}
		err = rows.Scan(
			&result.ID,
			&result.AddressID,
			&result.UserID,
			&result.IsDefault,
			&result.CreatedAt,
			&result.UpdatedAt,
			&result.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (a *addressUserRepository) Save(ctx context.Context, entity *models.AddressUser) (result *models.AddressUser, err error) {
	id, err := mysql_client.Save[models.AddressUser](ctx, a.dbx(ctx), queries.InsertAddressUser, &entity.AddressID, &entity.UserID)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (a *addressUserRepository) FindByID(ctx context.Context, id uint64) (result *models.AddressUser, err error) {
	result = &models.AddressUser{}
	err = mysql_client.FindByID[models.AddressUser](ctx, a.dbx(ctx), queries.SelectAddressUserByID, id, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *addressUserRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.AddressUser, err error) {
	return mysql_client.FindAll[models.AddressUser](ctx, a.dbx(ctx), queries.SelectAddressUsers)
}

var _ repository_interfaces.AddressUserRepository = &addressUserRepository{}

func newAddressUserRepository(dbx repository_interfaces.DBX) *addressUserRepository {
	return &addressUserRepository{
		dbx: dbx,
	}
}
