package domain

import (
	"context"
	"time"
)

//EmployeeFilter reqpresent query filter
type EmployeeFilter struct {
	IDs     []string
	Keyword string
	Num     int
	Cursor  string
	DeptIDs []string
}

// Employee represent employee data
type Employee struct {
	ID          string     `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	BirthPlace  string     `json:"birth_place"`
	DateOfBirth string     `json:"date_of_birth"`
	Title       string     `json:"title"`
	Department  Department `json:"department"`
	CreatedTime time.Time  `json:"created_time"`
	UpdatedTime time.Time  `json:"updated_time"`
}

// EmployeeRepository represent repository contract for employee
type EmployeeRepository interface {
	Create(ctx context.Context, e *Employee) (err error)
	Fetch(ctx context.Context, filter EmployeeFilter) (employees []Employee, nextCursor string, err error)
	Get(ctx context.Context, employeeID string) (employee Employee, err error)
	Update(ctx context.Context, e Employee) (employee Employee, err error)
	Delete(ctx context.Context, employeeID string) (err error)
}
