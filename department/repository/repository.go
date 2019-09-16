package repository

import (
	"context"
	"database/sql"

	employee "github.com/milhamhidayat/golang-clean-code-v2/domain"
)

// DepartmentRepository implement all method from interface
type DepartmentRepository struct {
	DB *sql.DB
}

// NewDepartmentRepository return new department repository
func NewDepartmentRepository(db *sql.DB) DepartmentRepository {
	return DepartmentRepository{
		DB: db,
	}
}

// Create is a repository to insert an article
func (r DepartmentRepository) Create(ctx context.Context, d employee.Department) (res employee.Department, err error) {
	return employee.Department{}, nil
}

// Fetch is a repository to fetch articles based on parameter
func (r DepartmentRepository) Fetch(ctx context.Context, filter employee.DepartmentFilter) (departments []employee.Department, nextCursor string, err error) {
	return []employee.Department{}, "", nil
}

// Get is a repository to get an article based on parameter
func (r DepartmentRepository) Get(ctx context.Context, departmentID string) (department employee.Department, err error) {
	return employee.Department{}, nil
}

// Update is a repository to update an article
func (r DepartmentRepository) Update(ctx context.Context, d employee.Department) (department employee.Department, err error) {
	return employee.Department{}, nil
}

// Delete is a repository to delete an article
func (r DepartmentRepository) Delete(ctx context.Context, departmentID string) (err error) {
	return nil
}
