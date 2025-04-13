package repositories

import repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"

type addressOwnerRepository struct {
	dbx repository_interfaces.DBX
}

func newAddressOwnerRepository(dbx repository_interfaces.DBX) *addressOwnerRepository {
	return &addressOwnerRepository{
		dbx: dbx,
	}
}
