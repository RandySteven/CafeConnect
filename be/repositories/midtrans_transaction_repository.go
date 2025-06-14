package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type midtransTransactionRepository struct {
	dbx repository_interfaces.DBX
}

func (m *midtransTransactionRepository) Save(ctx context.Context, entity *models.MidtransTransaction) (result *models.MidtransTransaction, err error) {
	id, err := mysql_client.Save[models.MidtransTransaction](ctx, m.dbx(ctx), queries.InsertMidtransTransaction, &entity.TransactionCode, &entity.TotalAmt, &entity.Token, &entity.RedirectURL)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (m *midtransTransactionRepository) FindByID(ctx context.Context, id uint64) (result *models.MidtransTransaction, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *midtransTransactionRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.MidtransTransaction, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *midtransTransactionRepository) FindByTransactionCode(ctx context.Context, transactionCode string) (result *models.MidtransTransaction, err error) {
	result = &models.MidtransTransaction{}
	err = m.dbx(ctx).QueryRowContext(ctx, queries.SelectMidtransByTransactionCode.String(), transactionCode).Scan(
		&result.ID,
		&result.TransactionCode,
		&result.TotalAmt,
		&result.Token,
		&result.RedirectURL,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

var _ repository_interfaces.MidtransTransactionRepository = &midtransTransactionRepository{}

func newMidtransTransactionRepository(dbx repository_interfaces.DBX) *midtransTransactionRepository {
	return &midtransTransactionRepository{
		dbx: dbx,
	}
}
