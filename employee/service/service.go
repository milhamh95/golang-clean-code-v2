package service

import (
	"context"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
)

// Service is an employee service
type Service struct {
	departmentRepo domain.DepartmentRepository
	employeeRepo   domain.EmployeeRepository
}

// New will crate a new employee service
func New(departmentRepo domain.DepartmentRepository, employeeRepo domain.EmployeeRepository) Service {
	return Service{
		departmentRepo: departmentRepo,
		employeeRepo:   employeeRepo,
	}
}

// Create will create a new employee
func (s Service) Create(ctx context.Context, e *domain.Employee) (err error) {
	return
}

// Fetch will return employess based on filter
func (s Service) Fetch(ctx context.Context, filter domain.EmployeeFilter) (employees []domain.Employee, nextCursor string, err error) {
	return
}

// Get will return an employee
func (s Service) Get(ctx context.Context, employeeID string) (employee domain.Employee, err error) {
	return
}

// Update will update an employee
func (s Service) Update(ctx context.Context, e domain.Employee) (employee domain.Employee, err error) {
	return
}

// Delete will delete an employee
func (s Service) Delete(ctx context.Context, employeeID string) (err error) {
	err = s.employeeRepo.Delete(context.Background(), employeeID)
	if err != nil {
		return
	}

	return
}
