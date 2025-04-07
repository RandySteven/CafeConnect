package repositories

import (
	"database/sql"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
)

type Repositories struct {
	UserRepository repository_interfaces.UserRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{}
}
