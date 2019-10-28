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
