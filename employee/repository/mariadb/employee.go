package mariadb

import (
	"context"
	"database/sql"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
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

// Create is a repository to insert an employee
func (r Repository) Create(ctx context.Context, e *domain.Employee) (err error) {
	return nil
}

// Get is a repository to get an employee
func (r Repository) Get(ctx context.Context, employeeID string) (employee domain.Employee, err error) {
	return
}

// Fetch is a repository to fetch employees
func (r Repository) Fetch(ctx context.Context, filter domain.EmployeeFilter) (employees []domain.Employee, nextCursor string, err error) {
	return
}

// Update is a repository to update an employee
func (r Repository) Update(ctx context.Context, e domain.Employee) (employee domain.Employee, err error) {
	return
}

// Delete is a repository to delete an employee
func (r Repository) Delete(ctx context.Context, employeeID string) (err error) {
	return
}
