package service_test

import (
	"context"
	"testing"

	"github.com/friendsofgo/errors"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	"github.com/milhamhidayat/golang-clean-code-v2/domain/mocks"
	"github.com/milhamhidayat/golang-clean-code-v2/employee/service"
	"github.com/milhamhidayat/golang-clean-code-v2/testdata"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	var employee domain.Employee
	testdata.UnmarshallGoldenToJSON(t, "employee-1S9XpJCvJbt1plvU36tAcJWS2ZW", &employee)

	mockDepartmentRepo := new(mocks.DepartmentRepository)
	mockEmployeeRepo := new(mocks.EmployeeRepository)

	tests := map[string]struct {
		employeeRepo map[string]testdata.FuncCall
		expectedErr  error
	}{
		"success": {
			employeeRepo: map[string]testdata.FuncCall{
				"Create": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), &employee},
					Output: []interface{}{nil},
				},
			},
			expectedErr: nil,
		},
		"with error create an employee": {
			employeeRepo: map[string]testdata.FuncCall{
				"Create": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), &employee},
					Output: []interface{}{errors.New("unexpected error")},
				},
			},
			expectedErr: errors.New("unexpected error"),
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			for name, fn := range tc.employeeRepo {
				if fn.Called {
					mockEmployeeRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			employeeService := service.New(mockDepartmentRepo, mockEmployeeRepo)
			err := employeeService.Create(context.Background(), &employee)

			mockEmployeeRepo.AssertExpectations(t)

			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestFetch(t *testing.T) {
	var employee1, employee2 domain.Employee
	testdata.UnmarshallGoldenToJSON(t, "employee-1S9XpJCvJbt1plvU36tAcJWS2ZW", &employee1)
	testdata.UnmarshallGoldenToJSON(t, "employee-1SYxHnSCbFCxLr7zUxk5j8cB0Cr", &employee2)

	mockDepartment := domain.Department{
		ID:   "1",
		Name: "Marketing",
	}

	mockDepartmentRepo := new(mocks.DepartmentRepository)
	mockEmployeeRepo := new(mocks.EmployeeRepository)

	tests := map[string]struct {
		filter         domain.EmployeeFilter
		employeeRepo   map[string]testdata.FuncCall
		departmentRepo map[string]testdata.FuncCall
		expectedRes    []domain.Employee
		expectedCursor string
		expectedErr    error
	}{
		"success with num": {
			filter: domain.EmployeeFilter{Num: 1},
			employeeRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.EmployeeFilter{Num: 1}},
					Output: []interface{}{[]domain.Employee{employee1}, "cursor-1", nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{mock.Anything, mock.AnythingOfType("string")},
					Output: []interface{}{mockDepartment, nil},
				},
			},
			expectedRes:    []domain.Employee{employee1},
			expectedCursor: "cursor-1",
			expectedErr:    nil,
		},
		"success with num and cursor": {
			filter: domain.EmployeeFilter{Num: 1, Cursor: "cursor-1"},
			employeeRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.EmployeeFilter{Num: 1, Cursor: "cursor-1"}},
					Output: []interface{}{[]domain.Employee{employee2}, "cursor-2", nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{mock.Anything, mock.AnythingOfType("string")},
					Output: []interface{}{mockDepartment, nil},
				},
			},
			expectedRes:    []domain.Employee{employee2},
			expectedCursor: "cursor-2",
			expectedErr:    nil,
		},
		"success with num and cursor end of page": {
			filter: domain.EmployeeFilter{Num: 1, Cursor: "cursor-2"},
			employeeRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.EmployeeFilter{Num: 1, Cursor: "cursor-2"}},
					Output: []interface{}{[]domain.Employee{}, "cursor-2", nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{Called: false},
			},
			expectedRes:    []domain.Employee{},
			expectedCursor: "cursor-2",
			expectedErr:    nil,
		},
		"success with ids": {
			filter: domain.EmployeeFilter{IDs: []string{employee2.ID}},
			employeeRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.EmployeeFilter{IDs: []string{employee2.ID}}},
					Output: []interface{}{[]domain.Employee{employee2}, "", nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{mock.Anything, mock.AnythingOfType("string")},
					Output: []interface{}{mockDepartment, nil},
				},
			},
			expectedRes:    []domain.Employee{employee2},
			expectedCursor: "",
			expectedErr:    nil,
		},
		"success with keyword": {
			filter: domain.EmployeeFilter{Keyword: "emilia"},
			employeeRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.EmployeeFilter{Keyword: "emilia"}},
					Output: []interface{}{[]domain.Employee{employee1}, "", nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{mock.Anything, mock.AnythingOfType("string")},
					Output: []interface{}{mockDepartment, nil},
				},
			},
			expectedRes:    []domain.Employee{employee1},
			expectedCursor: "",
			expectedErr:    nil,
		},
		"success with dept ids": {
			filter: domain.EmployeeFilter{DeptIDs: []string{"1"}},
			employeeRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.EmployeeFilter{DeptIDs: []string{"1"}}},
					Output: []interface{}{[]domain.Employee{employee1}, "", nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{mock.Anything, mock.AnythingOfType("string")},
					Output: []interface{}{mockDepartment, nil},
				},
			},
			expectedRes:    []domain.Employee{employee1},
			expectedCursor: "",
			expectedErr:    nil,
		},
		"error fetch employee repo": {
			filter: domain.EmployeeFilter{Num: 1},
			employeeRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.EmployeeFilter{Num: 1}},
					Output: []interface{}{[]domain.Employee{}, "", errors.New("unknown error")},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{Called: false},
			},
			expectedRes:    []domain.Employee{},
			expectedCursor: "",
			expectedErr:    errors.New("unknown error"),
		},
		"error get department": {
			filter: domain.EmployeeFilter{Num: 1},
			employeeRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.EmployeeFilter{Num: 1}},
					Output: []interface{}{[]domain.Employee{employee1}, "cursor-1", nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{mock.Anything, mock.AnythingOfType("string")},
					Output: []interface{}{domain.Department{}, errors.New("unknown error")},
				},
			},
			expectedRes:    []domain.Employee{},
			expectedCursor: "",
			expectedErr:    errors.New("unknown error"),
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			for name, fn := range tc.employeeRepo {
				if fn.Called {
					mockEmployeeRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			for name, fn := range tc.departmentRepo {
				if fn.Called {
					mockDepartmentRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			employeeService := service.New(mockDepartmentRepo, mockEmployeeRepo)
			res, nextCursor, err := employeeService.Fetch(context.Background(), tc.filter)

			mockEmployeeRepo.AssertExpectations(t)
			mockDepartmentRepo.AssertExpectations(t)

			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedCursor, nextCursor)
			require.Equal(t, tc.expectedRes, res)
		})
	}
}

func TestGet(t *testing.T) {
	var (
		employee   domain.Employee
		department domain.Department
	)
	testdata.UnmarshallGoldenToJSON(t, "employee-1S9XpJCvJbt1plvU36tAcJWS2ZW", &employee)
	testdata.UnmarshallGoldenToJSON(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &department)

	mockDepartmentRepo := new(mocks.DepartmentRepository)
	mockEmployeeRepo := new(mocks.EmployeeRepository)

	tests := map[string]struct {
		employeeRepo   map[string]testdata.FuncCall
		departmentRepo map[string]testdata.FuncCall
		expectedRes    domain.Employee
		expectedErr    error
	}{
		"success": {
			employeeRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), employee.ID},
					Output: []interface{}{employee, nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), employee.Department.ID},
					Output: []interface{}{department, nil},
				},
			},
			expectedRes: employee,
			expectedErr: nil,
		},
		"with error get employee": {
			employeeRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), employee.ID},
					Output: []interface{}{domain.Employee{}, errors.New("unknown error")},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{Called: false},
			},
			expectedRes: domain.Employee{},
			expectedErr: errors.New("unknown error"),
		},
		"with error get department": {
			employeeRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), employee.ID},
					Output: []interface{}{employee, nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), employee.Department.ID},
					Output: []interface{}{domain.Department{}, errors.New("unexpected error")},
				},
			},
			expectedRes: domain.Employee{},
			expectedErr: errors.New("unexpected error"),
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			for name, fn := range tc.employeeRepo {
				if fn.Called {
					mockEmployeeRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			for name, fn := range tc.departmentRepo {
				if fn.Called {
					mockDepartmentRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			employeeService := service.New(mockDepartmentRepo, mockEmployeeRepo)
			res, err := employeeService.Get(context.Background(), employee.ID)

			mockDepartmentRepo.AssertExpectations(t)
			mockEmployeeRepo.AssertExpectations(t)

			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, res, employee)
		})
	}
}

func TestUpdate(t *testing.T) {
	var (
		employee   domain.Employee
		department domain.Department
	)
	testdata.UnmarshallGoldenToJSON(t, "employee-1S9XpJCvJbt1plvU36tAcJWS2ZW", &employee)
	testdata.UnmarshallGoldenToJSON(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &department)

	mockDepartmentRepo := new(mocks.DepartmentRepository)
	mockEmployeeRepo := new(mocks.EmployeeRepository)

	newEmployee := employee
	newEmployee.LastName = "Diana"
	newEmployee.Department = department

	tests := map[string]struct {
		employeeRepo   map[string]testdata.FuncCall
		departmentRepo map[string]testdata.FuncCall
		expectedRes    domain.Employee
		expectedErr    error
	}{
		"success": {
			employeeRepo: map[string]testdata.FuncCall{
				"Update": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), newEmployee},
					Output: []interface{}{newEmployee, nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), newEmployee.Department.ID},
					Output: []interface{}{department, nil},
				},
			},
			expectedRes: newEmployee,
			expectedErr: nil,
		},
		"with error update an employee": {
			employeeRepo: map[string]testdata.FuncCall{
				"Update": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), newEmployee},
					Output: []interface{}{domain.Employee{}, errors.New("unexpected error")},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), newEmployee.Department.ID},
					Output: []interface{}{department, nil},
				},
			},
			expectedRes: domain.Employee{},
			expectedErr: errors.New("unexpected error"),
		},
		"with error get a department": {
			employeeRepo: map[string]testdata.FuncCall{
				"Update": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), newEmployee},
					Output: []interface{}{newEmployee, nil},
				},
			},
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), newEmployee.Department.ID},
					Output: []interface{}{domain.Department{}, errors.New("unknown error")},
				},
			},
			expectedRes: domain.Employee{},
			expectedErr: errors.New("unknown error"),
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			for name, fn := range tc.employeeRepo {
				if fn.Called {
					mockEmployeeRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			for name, fn := range tc.departmentRepo {
				if fn.Called {
					mockDepartmentRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			employeeService := service.New(mockDepartmentRepo, mockEmployeeRepo)
			res, err := employeeService.Update(context.Background(), newEmployee)

			mockEmployeeRepo.AssertExpectations(t)
			mockDepartmentRepo.AssertExpectations(t)

			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, newEmployee, res)
		})
	}
}

func TestDelete(t *testing.T) {
	var employee domain.Employee
	testdata.UnmarshallGoldenToJSON(t, "employee-1S9XpJCvJbt1plvU36tAcJWS2ZW", &employee)

	mockDepartmentRepo := new(mocks.DepartmentRepository)
	mockEmployeeRepo := new(mocks.EmployeeRepository)

	tests := map[string]struct {
		employeeRepo map[string]testdata.FuncCall
		expectedErr  error
	}{
		"success": {
			employeeRepo: map[string]testdata.FuncCall{
				"Delete": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), employee.ID},
					Output: []interface{}{nil},
				},
			},
			expectedErr: nil,
		},
		"with error from employee repo": {
			employeeRepo: map[string]testdata.FuncCall{
				"Delete": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), employee.ID},
					Output: []interface{}{errors.New("unexpected error")},
				},
			},
			expectedErr: errors.New("unexpected error"),
		},
		"not found": {
			employeeRepo: map[string]testdata.FuncCall{
				"Delete": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), employee.ID},
					Output: []interface{}{domain.ErrNotFound},
				},
			},
			expectedErr: domain.ErrNotFound,
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			for name, fn := range tc.employeeRepo {
				if fn.Called {
					mockEmployeeRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			employeeService := service.New(mockDepartmentRepo, mockEmployeeRepo)
			err := employeeService.Delete(context.Background(), employee.ID)

			mockEmployeeRepo.AssertExpectations(t)

			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}
