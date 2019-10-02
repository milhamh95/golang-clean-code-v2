package mariadb

import (
	"database/sql"
)

// Repository implement all employee repository method from interface
type Repository struct {
	DB *sql.DB
}

// New return new department repository
func New(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}
