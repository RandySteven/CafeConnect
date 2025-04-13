package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
)

type addressRepository struct {
	dbx repository_interfaces.DBX
}

func (a addressRepository) Save(ctx context.Context, entity *models.Address) (result *models.Address, err error) {
	return
}

func (a addressRepository) FindByID(ctx context.Context, id uint64) (result *models.Address, err error) {
	return
}

func (a addressRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.Address, err error) {
	return
}

func (a addressRepository) Update(ctx context.Context, entity *models.Address) (result *models.Address, err error) {
	return
}

var _ repository_interfaces.AddressRepository = &addressRepository{}

func newAddressRepository(dbx repository_interfaces.DBX) *addressRepository {
	return &addressRepository{
		dbx: dbx,
	}
}
