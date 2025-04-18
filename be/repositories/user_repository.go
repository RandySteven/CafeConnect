package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
	"log"
)

const (
	userName int = iota + 1
	phoneNumber
	email
)

type userRepository struct {
	dbx repository_interfaces.DBX
}

func (u *userRepository) Save(ctx context.Context, entity *models.User) (result *models.User, err error) {
	id, err := mysql_client.Save[models.User](ctx, u.dbx(ctx), queries.InsertUser,
		&entity.Name, &entity.Username, &entity.Email, &entity.Password, &entity.PhoneNumber,
		&entity.ProfilePicture, &entity.DoB)
	if err != nil {
		log.Println("err query ", err)
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (u *userRepository) FindByID(ctx context.Context, id uint64) (result *models.User, err error) {
	result = &models.User{}
	err = mysql_client.FindByID[models.User](ctx, u.dbx(ctx), queries.SelectUserByID, id, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.User, err error) {
	return mysql_client.FindAll[models.User](ctx, u.dbx(ctx), queries.SelectUsers)
}

func (u *userRepository) Update(ctx context.Context, entity *models.User) (result *models.User, err error) {
	err = mysql_client.Update[models.User](ctx, u.dbx(ctx), ``, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (u *userRepository) DeleteByID(ctx context.Context, id uint64) (err error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindByEmail(ctx context.Context, emailStr string) (result *models.User, err error) {
	return u.finderBy(ctx, email, emailStr)
}

func (u *userRepository) FindByUsername(ctx context.Context, usernameStr string) (result *models.User, err error) {
	return u.finderBy(ctx, userName, usernameStr)
}

func (u *userRepository) FindByPhoneNumber(ctx context.Context, phoneNumberStr string) (result *models.User, err error) {
	return u.finderBy(ctx, phoneNumber, phoneNumberStr)
}

func (u *userRepository) finderBy(ctx context.Context, identify int, identifier string) (result *models.User, err error) {
	var query queries.GoQuery = ``
	switch identify {
	case userName:
		query = queries.SelectUsername
		break
	case phoneNumber:
		query = queries.SelectPhoneNumber
		break
	case email:
		query = queries.SelectEmail
		break
	}
	result = &models.User{}
	err = u.dbx(ctx).QueryRowContext(ctx, query.String(), identifier).Scan(
		&result.ID, &result.Name, &result.Email, &result.Username, &result.Password, &result.PhoneNumber,
		&result.ProfilePicture, &result.DoB, &result.CreatedAt, &result.UpdatedAt, &result.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

var _ repository_interfaces.UserRepository = &userRepository{}

func newUserRepository(dbx repository_interfaces.DBX) *userRepository {
	return &userRepository{
		dbx: dbx,
	}
}
