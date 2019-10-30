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
		err = errors.Wrap(err, "failed to create a department")
		return
	}

	return
}

// Fetch is a service to fetch department
func (s Service) Fetch(ctx context.Context, filter domain.DepartmentFilter) (departments []domain.Department, nextCursor string, err error) {
	departments, nextCursor, err = s.Repository.Fetch(ctx, filter)
	if err != nil {
		nextCursor = filter.Cursor
		err = errors.Wrap(err, "failed to fetch departments")
		return
	}

	return
}

// Get is a service to get a department
func (s Service) Get(ctx context.Context, departmentID string) (department domain.Department, err error) {
	department, err = s.Repository.Get(ctx, departmentID)
	if err != nil {
		err = errors.Wrap(err, "failed to get a department")
		return
	}

	return
}

// Update is a service to update a department
func (s Service) Update(ctx context.Context, d domain.Department) (department domain.Department, err error) {
	department, err = s.Repository.Update(ctx, d)
	if err != nil {
		err = errors.Wrap(err, "failed to update a department")
		return
	}
	return
}

// Delete is a service to delete a department
func (s Service) Delete(ctx context.Context, departmentID string) (err error) {
	err = s.Repository.Delete(context.Background(), departmentID)
	if err != nil {
		err = errors.Wrap(err, "failed to delete a department")
		return
	}

	return
}
