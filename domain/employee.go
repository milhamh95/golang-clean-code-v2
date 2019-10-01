package employee

type EmployeeFilter struct {
	IDs     []string
	Keyword string
	Num     int
	Cursor  string
	DeptIDs []string
}