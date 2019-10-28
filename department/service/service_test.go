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
	tests := map[string]struct {
		departmentRepo map[string]testdata.FuncCall
		expectedError  error
	}{
		"success": {
			departmentRepo: map[string]testdata.FuncCall{
				"Create": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), &domain.Department{}},
					Output: []interface{}{nil},
				},
			},
			expectedError: nil,
		},
		"error": {
			departmentRepo: map[string]testdata.FuncCall{
				"Create": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), &domain.Department{}},
					Output: []interface{}{errors.New("unexpected error")},
				},
			},
			expectedError: errors.New("unexpected error"),
		},
	}

	mockDepartmentRepo := new(mocks.DepartmentRepository)

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			for name, fn := range tc.departmentRepo {
				if fn.Called {
					mockDepartmentRepo.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}

			departmentService := service.New(mockDepartmentRepo)
			err := departmentService.Create(context.Background(), &domain.Department{})

			mockDepartmentRepo.AssertExpectations(t)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)

		})
	}
}
