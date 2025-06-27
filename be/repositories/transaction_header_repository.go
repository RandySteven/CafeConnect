package repositories

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	mysql_client "github.com/RandySteven/CafeConnect/be/pkg/mysql"
	"github.com/RandySteven/CafeConnect/be/queries"
)

type transactionHeaderRepository struct {
	dbx repository_interfaces.DBX
}

func (t *transactionHeaderRepository) Update(ctx context.Context, entity *models.TransactionHeader) (result *models.TransactionHeader, err error) {
	err = mysql_client.Update[models.TransactionHeader](ctx, t.dbx(ctx), queries.UpdateTransactionHeader,
		&entity.UserID,
		&entity.CafeID,
		&entity.TransactionCode,
		&entity.Status,
		&entity.TransactionAt,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.ID,
	)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (t *transactionHeaderRepository) CreateIndex(ctx context.Context) (err error) {
	_, err = t.dbx(ctx).ExecContext(ctx, queries.TransactionCodeIndex.String())
	if err != nil {
		return err
	}
	return nil
}

func (t *transactionHeaderRepository) DropIndex(ctx context.Context) (err error) {
	_, err = t.dbx(ctx).ExecContext(ctx, fmt.Sprintf(queries.DropIndex.String(), `transaction_code`, `transactions`))
	if err != nil {
		return err
	}
	return nil
}

func (t *transactionHeaderRepository) Save(ctx context.Context, entity *models.TransactionHeader) (result *models.TransactionHeader, err error) {
	id, err := mysql_client.Save[models.TransactionHeader](ctx, t.dbx(ctx), queries.InsertTransactionHeader, &entity.UserID, &entity.CafeID, &entity.TransactionCode, &entity.Status, &entity.TransactionAt)
	if err != nil {
		return nil, err
	}
	entity.ID = *id
	return entity, nil
}

func (t *transactionHeaderRepository) FindByTransactionCode(ctx context.Context, transactionCode string) (result *models.TransactionHeader, err error) {
	result = &models.TransactionHeader{}
	err = t.dbx(ctx).QueryRowContext(ctx, queries.SelectTransactionHeaderByTransactionCode.String(), transactionCode).
		Scan(
			&result.ID,
			&result.UserID,
			&result.CafeID,
			&result.TransactionCode,
			&result.Status,
			&result.TransactionAt,
			&result.CreatedAt,
			&result.UpdatedAt,
			&result.DeletedAt,
		)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *transactionHeaderRepository) FindByUserID(ctx context.Context, userId uint64) (result []*models.TransactionHeader, err error) {
	rows, err := t.dbx(ctx).QueryContext(ctx, queries.SelectTransactionHeaderByUserID.String(), userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		header := &models.TransactionHeader{}
		err = rows.Scan(
			&header.ID,
			&header.UserID,
			&header.CafeID,
			&header.TransactionCode,
			&header.Status,
			&header.TransactionAt,
			&header.CreatedAt,
			&header.UpdatedAt,
			&header.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, header)
	}
	return result, nil
}

func (t *transactionHeaderRepository) FindByTransactionStatus(ctx context.Context, status string) (result []*models.TransactionHeader, err error) {
	rows, err := t.dbx(ctx).QueryContext(ctx, queries.SelectTransactionHeadersByStatus.String(), status)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		header := &models.TransactionHeader{}
		err = rows.Scan(
			&header.ID,
			&header.UserID,
			&header.CafeID,
			&header.TransactionCode,
			&header.Status,
			&header.TransactionAt,
			&header.CreatedAt,
			&header.UpdatedAt,
			&header.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, header)
	}
	return result, nil
}

var _ repository_interfaces.TransactionHeaderRepository = &transactionHeaderRepository{}

func newTransactionHeaderRepository(dbx repository_interfaces.DBX) *transactionHeaderRepository {
	return &transactionHeaderRepository{
		dbx: dbx,
	}
}
