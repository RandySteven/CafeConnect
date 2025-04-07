package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
)

type userRepository struct {
	dbx repository_interfaces.DBX
}

func (u *userRepository) Save(ctx context.Context, entity *models.User) (result *models.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindByID(ctx context.Context, id uint64) (result *models.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) Update(ctx context.Context, entity *models.User) (result *models.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) DeleteByID(ctx context.Context, id uint64) (err error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (result *models.User, err error) {
	//TODO implement me
	panic("implement me")
}

var _ repository_interfaces.UserRepository = &userRepository{}

func newUserRepository(dbx repository_interfaces.DBX) *userRepository {
	return &userRepository{
		dbx: dbx,
	}
}
