package repositories

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type transactionDetailRepository struct {
	dbx repository_interfaces.DBX
}

func (t *transactionDetailRepository) FindByTransactionId(ctx context.Context, transactionId uint64) (result []*models.TransactionDetail, err error) {
	rows, err := t.dbx(ctx).QueryContext(ctx, queries.SelectTransactionDetailsByTransactionId.String(), transactionId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		transactionDetail := &models.TransactionDetail{}
		err = rows.Scan(
			&transactionDetail.ID,
			&transactionDetail.TransactionID,
			&transactionDetail.CafeProductID,
			&transactionDetail.Qty,
			&transactionDetail.CreatedAt,
			&transactionDetail.UpdatedAt,
			&transactionDetail.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, transactionDetail)
	}

	return result, nil
}

func (t *transactionDetailRepository) Save(ctx context.Context, entity *models.TransactionDetail) (result *models.TransactionDetail, err error) {
	id, err := mysql_client.Save[models.TransactionDetail](ctx, t.dbx(ctx), queries.InsertTransactionDetail, &entity.TransactionID, &entity.CafeProductID, &entity.Qty)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (t *transactionDetailRepository) FindByID(ctx context.Context, id uint64) (result *models.TransactionDetail, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *transactionDetailRepository) FindAll(ctx context.Context, skip uint64, take uint64) (result []*models.TransactionDetail, err error) {
	//TODO implement me
	panic("implement me")
}

var _ repository_interfaces.TransactionDetailRepository = &transactionDetailRepository{}

func newTransactionDetailRepository(dbx repository_interfaces.DBX) *transactionDetailRepository {
	return &transactionDetailRepository{
		dbx: dbx,
	}
}
