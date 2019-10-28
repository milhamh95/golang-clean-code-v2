package service

import (
	"context"

	"github.com/friendsofgo/errors"

	domain "github.com/milhamhidayat/golang-clean-code-v2/domain"
)

// Service is a department service
type Service struct {
	Repository domain.DepartmentRepository
}

// New will return a department service
func New(
	repo domain.DepartmentRepository,
) domain.DepartmentRepository {
	return Service{
		Repository: repo,
	}
}

// Create is a service to create department
func (s Service) Create(ctx context.Context, d *domain.Department) (err error) {
	err = s.Repository.Create(ctx, d)
	if err != nil {
		return
	}

	return
}

// Fetch is a service to fetch department
func (s Service) Fetch(ctx context.Context, filter domain.DepartmentFilter) (departments []domain.Department, nextCursor string, err error) {
	return []domain.Department{}, "", errors.New("not yet implemented")
}

// Get is a service to get a department
func (s Service) Get(ctx context.Context, departmentID string) (department domain.Department, err error) {
	department, err = s.Repository.Get(ctx, departmentID)
	if err != nil {
		return
	}

	return
}

// Update is a service to update a department
func (s Service) Update(ctx context.Context, d domain.Department) (department domain.Department, err error) {
	return domain.Department{}, errors.New("not yet implemented")
}

// Delete is a service to delete a department
func (s Service) Delete(ctx context.Context, departmentID string) (err error) {
	err = s.Repository.Delete(context.Background(), departmentID)
	if err != nil {
		return
	}

	return
}
