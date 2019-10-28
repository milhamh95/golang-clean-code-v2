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
