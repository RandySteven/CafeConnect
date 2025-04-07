package repository_interfaces

import (
	"context"
	"database/sql"
	"github.com/RandySteven/CafeConnect/be/apperror"
)

type (
	Many2ManyRepository[T any] interface {
		Save(ctx context.Context, entity []*T) (result *T, err error)
	}

	Saver[T any] interface {
		Save(ctx context.Context, entity *T) (result *T, err error)
	}

	Finder[T any] interface {
		FindByID(ctx context.Context, id uint64) (result *T, err error)
		FindAll(ctx context.Context, skip uint64, take uint64) (result []*T, err error)
	}

	Updater[T any] interface {
		Update(ctx context.Context, entity *T) (result *T, err error)
	}

	Deleter[T any] interface {
		DeleteByID(ctx context.Context, id uint64) (err error)
	}

	DBX func(ctx context.Context) Trigger

	Trigger interface {
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}

	Transaction interface {
		RunInTx(ctx context.Context, txFunc func(ctx context.Context) (customErr *apperror.CustomError)) (customErr *apperror.CustomError)
	}
)
