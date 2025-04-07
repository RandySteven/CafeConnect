package repositories

import (
	"context"
	"database/sql"
	"github.com/RandySteven/CafeConnect/be/apperror"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
)

type (
	transaction struct {
		db *sql.DB
	}

	txCtxKey struct{}
)

func (t *transaction) RunInTx(ctx context.Context, txFunc func(ctx context.Context) *apperror.CustomError) (customErr *apperror.CustomError) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.NewCustomError(apperror.ErrInternalServer, `failed to begin tx`, err)
	}
	txCtx := txToContext(ctx, tx)
	if customErr = txFunc(txCtx); customErr != nil {
		_ = tx.Rollback()
		return customErr
	}

	if err = tx.Commit(); err != nil {
		return apperror.NewCustomError(apperror.ErrInternalServer, `failed to commit tx`, err)
	}

	return nil
}

func txToContext(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}

func txFromContext(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(txCtxKey{}).(*sql.Tx)
	if ok {
		return tx
	}

	return nil
}

var _ repository_interfaces.Transaction = &transaction{}

func newTransaction(db *sql.DB) (*transaction, repository_interfaces.DBX) {
	return &transaction{
			db: db,
		}, func(ctx context.Context) repository_interfaces.Trigger {
			if tx := txFromContext(ctx); tx != nil {
				return tx
			}
			return db
		}
}
