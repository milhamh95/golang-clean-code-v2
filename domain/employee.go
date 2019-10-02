package domain

import "time"

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
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	BirthPlace  string    `json:"birth_place"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Title       string    `json:"title"`
	Department  Department
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}
