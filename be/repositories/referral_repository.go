package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type referralRepository struct {
	dbx repository_interfaces.DBX
}

func (r *referralRepository) Update(ctx context.Context, entity *models.Referral) (result *models.Referral, err error) {
	err = mysql_client.Update[models.Referral](ctx, r.dbx(ctx), queries.UpdateReferralByID, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *referralRepository) FindByUserID(ctx context.Context, userId uint64) (result *models.Referral, err error) {
	result = &models.Referral{}
	err = r.dbx(ctx).QueryRowContext(ctx, queries.SelectReferralByUserID.String(), userId).Scan(
		&result.ID,
		&result.Code,
		&result.UserID,
		&result.ExpiredTime,
		&result.Status,
		&result.NumbOfUsage,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *referralRepository) FindByCode(ctx context.Context, code string) (result *models.Referral, err error) {
	result = &models.Referral{}
	err = r.dbx(ctx).QueryRowContext(ctx, queries.SelectReferralByCode.String(), code).Scan(
		&result.ID,
		&result.Code,
		&result.UserID,
		&result.ExpiredTime,
		&result.Status,
		&result.NumbOfUsage,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *referralRepository) Save(ctx context.Context, entity *models.Referral) (result *models.Referral, err error) {
	id, err := mysql_client.Save[models.Referral](ctx, r.dbx(ctx), queries.InsertReferral, &entity.Code, &entity.UserID, &entity.ExpiredTime, &entity.Status)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (r *referralRepository) FindByID(ctx context.Context, id uint64) (result *models.Referral, err error) {
	result = &models.Referral{}
	err = mysql_client.FindByID[models.Referral](ctx, r.dbx(ctx), queries.SelectReferralByID, id, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *referralRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.Referral, err error) {
	//TODO implement me
	panic("implement me")
}

var _ repository_interfaces.ReferralRepository = &referralRepository{}

func newReferralRepository(dbx repository_interfaces.DBX) *referralRepository {
	return &referralRepository{
		dbx: dbx,
	}
}
