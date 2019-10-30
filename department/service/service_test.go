package service_test

import (
	"context"
	"testing"

	"github.com/friendsofgo/errors"
	"github.com/stretchr/testify/require"

	"github.com/milhamhidayat/golang-clean-code-v2/department/service"
	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	"github.com/milhamhidayat/golang-clean-code-v2/domain/mocks"
	"github.com/milhamhidayat/golang-clean-code-v2/testdata"
)

func TestCreate(t *testing.T) {
	var department domain.Department
	testdata.UnmarshallGoldenToJSON(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &department)

	mockDepartmentRepo := new(mocks.DepartmentRepository)

	tests := map[string]struct {
		departmentRepo map[string]testdata.FuncCall
		expectedError  error
	}{
		"success": {
			departmentRepo: map[string]testdata.FuncCall{
				"Create": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), &department},
					Output: []interface{}{nil},
				},
			},
			expectedError: nil,
		},
		"error": {
			departmentRepo: map[string]testdata.FuncCall{
				"Create": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), &department},
					Output: []interface{}{errors.New("unexpected error")},
				},
			},
			expectedError: errors.New("unexpected error"),
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			for name, fn := range tc.departmentRepo {
				if fn.Called {
					mockDepartmentRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			departmentService := service.New(mockDepartmentRepo)
			err := departmentService.Create(context.Background(), &department)

			mockDepartmentRepo.AssertExpectations(t)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestFetch(t *testing.T) {
	var department1, department2, department3 domain.Department
	testdata.UnmarshallGoldenToJSON(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &department1)
	testdata.UnmarshallGoldenToJSON(t, "department-0ujssxh0cECutqzMgbtXSGnjorm", &department2)
	testdata.UnmarshallGoldenToJSON(t, "department-0ujsszgFvbiEr7CDgE3z8MAUPFt", &department3)

	mockDepartmentRepo := new(mocks.DepartmentRepository)

	tests := map[string]struct {
		filter         domain.DepartmentFilter
		departmentRepo map[string]testdata.FuncCall
		expectedRes    []domain.Department
		expectedCursor string
		expectedErr    error
	}{
		"success with num first page": {
			filter: domain.DepartmentFilter{Num: 2},
			departmentRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.DepartmentFilter{Num: 2}},
					Output: []interface{}{[]domain.Department{department1, department2}, "MHVqc3N4aDBjRUN1dHF6TWdidFhTR25qb3Jt", nil},
				},
			},
			expectedRes:    []domain.Department{department1, department2},
			expectedCursor: "MHVqc3N4aDBjRUN1dHF6TWdidFhTR25qb3Jt",
			expectedErr:    nil,
		},
		"success with num and cursor second page": {
			filter: domain.DepartmentFilter{Num: 2, Cursor: "MHVqc3N4aDBjRUN1dHF6TWdidFhTR25qb3Jt"},
			departmentRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.DepartmentFilter{Num: 2, Cursor: "MHVqc3N4aDBjRUN1dHF6TWdidFhTR25qb3Jt"}},
					Output: []interface{}{[]domain.Department{department3}, "MHVqc3N6Z0Z2YmlFcjdDRGdFM3o4TUFVUEZ0", nil},
				},
			},
			expectedRes:    []domain.Department{department3},
			expectedCursor: "MHVqc3N6Z0Z2YmlFcjdDRGdFM3o4TUFVUEZ0",
			expectedErr:    nil,
		},
		"success with num and cursor end of page": {
			filter: domain.DepartmentFilter{Num: 2, Cursor: "MHVqc3N6Z0Z2YmlFcjdDRGdFM3o4TUFVUEZ0"},
			departmentRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.DepartmentFilter{Num: 2, Cursor: "MHVqc3N6Z0Z2YmlFcjdDRGdFM3o4TUFVUEZ0"}},
					Output: []interface{}{[]domain.Department{}, "MHVqc3N6Z0Z2YmlFcjdDRGdFM3o4TUFVUEZ0", nil},
				},
			},
			expectedRes:    []domain.Department{},
			expectedCursor: "MHVqc3N6Z0Z2YmlFcjdDRGdFM3o4TUFVUEZ0",
			expectedErr:    nil,
		},
		"succes with ids": {
			filter: domain.DepartmentFilter{IDs: []string{department1.ID, department3.ID}},
			departmentRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.DepartmentFilter{IDs: []string{department1.ID, department3.ID}}},
					Output: []interface{}{[]domain.Department{department1, department3}, "", nil},
				},
			},
			expectedRes:    []domain.Department{department1, department3},
			expectedCursor: "",
			expectedErr:    nil,
		},
		"success with keyword": {
			filter: domain.DepartmentFilter{Keyword: "marketing"},
			departmentRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.DepartmentFilter{Keyword: "marketing"}},
					Output: []interface{}{[]domain.Department{department1}, "", nil},
				},
			},
			expectedRes:    []domain.Department{department1},
			expectedCursor: "",
			expectedErr:    nil,
		},
		"error fetch department": {
			filter: domain.DepartmentFilter{Num: 2},
			departmentRepo: map[string]testdata.FuncCall{
				"Fetch": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), domain.DepartmentFilter{Num: 2}},
					Output: []interface{}{[]domain.Department{}, "", errors.New("unknown error")},
				},
			},
			expectedRes:    []domain.Department{},
			expectedCursor: "",
			expectedErr:    errors.New("unknown error"),
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			for name, fn := range tc.departmentRepo {
				if fn.Called {
					mockDepartmentRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			departmentService := service.New(mockDepartmentRepo)
			res, cursor, err := departmentService.Fetch(context.Background(), tc.filter)

			mockDepartmentRepo.AssertExpectations(t)

			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, res, tc.expectedRes)
			require.Equal(t, cursor, tc.expectedCursor)
		})
	}
}

func TestGet(t *testing.T) {
	var department domain.Department
	testdata.UnmarshallGoldenToJSON(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &department)

	mockDepartmentRepo := new(mocks.DepartmentRepository)

	tests := map[string]struct {
		departmentRepo map[string]testdata.FuncCall
		expectedRes    domain.Department
		expectedErr    error
	}{
		"success": {
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), department.ID},
					Output: []interface{}{department, nil},
				},
			},
			expectedRes: department,
			expectedErr: nil,
		},
		"success with department not found": {
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), department.ID},
					Output: []interface{}{domain.Department{}, errors.New("department is not found")},
				},
			},
			expectedRes: domain.Department{},
			expectedErr: errors.New("department is not found"),
		},
		"with error from department repository": {
			departmentRepo: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), department.ID},
					Output: []interface{}{domain.Department{}, errors.New("unexpected error")},
				},
			},
			expectedRes: domain.Department{},
			expectedErr: errors.New("unexpected error"),
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			for name, fn := range tc.departmentRepo {
				if fn.Called {
					mockDepartmentRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			departmentService := service.New(mockDepartmentRepo)
			res, err := departmentService.Get(context.Background(), department.ID)

			mockDepartmentRepo.AssertExpectations(t)

			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
				return
			}

			require.Equal(t, department, res)
			require.NoError(t, err)
		})
	}
}

func TestDelete(t *testing.T) {
	var department domain.Department
	testdata.UnmarshallGoldenToJSON(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &department)

	mockDepartmentRepo := new(mocks.DepartmentRepository)

	tests := map[string]struct {
		departmentRepo map[string]testdata.FuncCall
		expectedErr    error
	}{
		"success": {
			departmentRepo: map[string]testdata.FuncCall{
				"Delete": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), department.ID},
					Output: []interface{}{nil},
				},
			},
			expectedErr: nil,
		},
		"with error from department repo": {
			departmentRepo: map[string]testdata.FuncCall{
				"Delete": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), department.ID},
					Output: []interface{}{errors.New("unexpected error")},
				},
			},
			expectedErr: errors.New("unexpected error"),
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			for name, fn := range tc.departmentRepo {
				if fn.Called {
					mockDepartmentRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			departmentService := service.New(mockDepartmentRepo)
			err := departmentService.Delete(context.Background(), department.ID)

			mockDepartmentRepo.AssertExpectations(t)

			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}
