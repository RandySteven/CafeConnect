package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type verifyTokenRepository struct {
	dbx repository_interfaces.DBX
}

func (v *verifyTokenRepository) Save(ctx context.Context, entity *models.VerifyToken) (result *models.VerifyToken, err error) {
	id, err := mysql_client.Save[models.VerifyToken](ctx, v.dbx(ctx), queries.InsertVerifyToken,
		&entity.Token, &entity.UserID, &entity.IsClicked, &entity.ExpiredTime)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (v *verifyTokenRepository) FindByToken(ctx context.Context, token string) (result *models.VerifyToken, err error) {
	result = &models.VerifyToken{}
	err = v.dbx(ctx).QueryRowContext(ctx, queries.SelectVerifyTokenByToken.String(), token).Scan(
		&result.ID, &result.Token, &result.UserID, &result.IsClicked, &result.ExpiredTime, &result.CreatedAt, &result.UpdatedAt, &result.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *verifyTokenRepository) Update(ctx context.Context, entity *models.VerifyToken) (result *models.VerifyToken, err error) {
	err = mysql_client.Update[models.VerifyToken](ctx, v.dbx(ctx), queries.UpdateVerifyToken,
		&entity.Token, &entity.UserID, &entity.IsClicked, &entity.ExpiredTime, &entity.CreatedAt, &entity.UpdatedAt, &entity.DeletedAt, &entity.ID,
	)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

var _ repository_interfaces.VerifyTokenRepository = &verifyTokenRepository{}

func newVerifyTokenRepository(dbx repository_interfaces.DBX) *verifyTokenRepository {
	return &verifyTokenRepository{
		dbx: dbx,
	}
}
