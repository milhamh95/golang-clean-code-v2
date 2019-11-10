package service

import (
	"context"

	"golang.org/x/sync/errgroup"

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
	err = s.employeeRepo.Create(ctx, e)
	if err != nil {
		return
	}

	return
}

// Fetch will return employess based on filter
func (s Service) Fetch(ctx context.Context, filter domain.EmployeeFilter) (employees []domain.Employee, nextCursor string, err error) {
	employees, nextCursor, err = s.employeeRepo.Fetch(ctx, filter)
	if err != nil {
		return
	}

	if len(employees) == 0 {
		return
	}

	err = s.fetchDepartment(ctx, employees)
	if err != nil {
		return
	}

	return
}

func (s Service) fetchDepartment(ctx context.Context, e []domain.Employee) (err error) {
	empDept := map[string]domain.Department{}
	for _, v := range e {
		empDept[v.ID] = domain.Department{}
	}

	g, ctx := errgroup.WithContext(context.Background())
	c := make(chan domain.Department)
	for k := range empDept {
		k := k
		g.Go(func() error {
			dept, err := s.departmentRepo.Get(ctx, k)
			if err != nil {
				return err
			}
			c <- dept
			return nil
		})
	}

	go func() {
		g.Wait()
		close(c)
	}()

	for v := range c {
		if v == (domain.Department{}) {
			empDept[v.ID] = v
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	for i, v := range e {
		if d, ok := empDept[v.Department.ID]; ok {
			e[i].Department = d
		}
	}

	return
}

// Get will return an employee
func (s Service) Get(ctx context.Context, employeeID string) (employee domain.Employee, err error) {
	employee, err = s.employeeRepo.Get(ctx, employeeID)
	if err != nil {
		return
	}

	department, err := s.departmentRepo.Get(ctx, employee.Department.ID)
	if err != nil {
		return
	}

	employee.Department = department

	return
}

// Update will update an employee
func (s Service) Update(ctx context.Context, e domain.Employee) (employee domain.Employee, err error) {
	ch1 := make(chan func() (domain.Employee, error))
	ch2 := make(chan func() (domain.Department, error))

	go func(ch chan func() (domain.Employee, error), e domain.Employee) {
		employee, err := s.employeeRepo.Update(ctx, e)

		ch <- (func() (domain.Employee, error) {
			return employee, err
		})
	}(ch1, e)

	go func(ch chan func() (domain.Department, error), id string) {
		department, err := s.departmentRepo.Get(ctx, id)

		ch <- (func() (domain.Department, error) {
			return department, err
		})
	}(ch2, e.Department.ID)

	employee, err = (<-ch1)()
	if err != nil {
		return
	}

	department, err := (<-ch2)()
	if err != nil {
		return
	}

	employee.Department = department

	return
}

// Delete will delete an employee
func (s Service) Delete(ctx context.Context, employeeID string) (err error) {
	err = s.employeeRepo.Delete(ctx, employeeID)
	if err != nil {
		return
	}

	return
}
