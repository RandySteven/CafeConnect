package repositories

import (
	"database/sql"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
)

type Repositories struct {
	UserRepository  repository_interfaces.UserRepository
	PointRepository repository_interfaces.PointRepository
	Transaction     repository_interfaces.Transaction
}

func NewRepositories(db *sql.DB) *Repositories {
	transaction, dbx := newTransaction(db)
	return &Repositories{
		UserRepository:  newUserRepository(dbx),
		PointRepository: newPointRepository(dbx),
		Transaction:     transaction,
	}
}
