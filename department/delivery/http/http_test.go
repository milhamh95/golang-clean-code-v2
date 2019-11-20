package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/friendsofgo/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	handler "github.com/milhamhidayat/golang-clean-code-v2/department/delivery/http"
	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	"github.com/milhamhidayat/golang-clean-code-v2/domain/mocks"
	"github.com/milhamhidayat/golang-clean-code-v2/pkg/middleware"
	"github.com/milhamhidayat/golang-clean-code-v2/testdata"
)

func TestInsert(t *testing.T) {
	e := testdata.GetEchoServer()

	var mockDepartment domain.Department
	rawMockDepartment := testdata.GetGolden(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K")
	err := json.Unmarshal(rawMockDepartment, &mockDepartment)
	require.NoError(t, err)

	mockDepartmentService := new(mocks.DepartmentService)

	tests := map[string]struct {
		reqBody           []byte
		departmentService map[string]testdata.FuncCall
		expectedStatus    int
	}{
		"success": {
			reqBody: rawMockDepartment,
			departmentService: map[string]testdata.FuncCall{
				"Create": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), &mockDepartment},
					Output: []interface{}{nil},
				},
			},
			expectedStatus: http.StatusCreated,
		},
		"invalid request body": {
			reqBody: []byte(``),
			departmentService: map[string]testdata.FuncCall{
				"Create": testdata.FuncCall{Called: false},
			},
			expectedStatus: http.StatusBadRequest,
		},
		"missing department name attribute": {
			reqBody: []byte(`
				{
					"description": "new department"
				}
			`),
			departmentService: map[string]testdata.FuncCall{
				"Create": testdata.FuncCall{Called: false},
			},
			expectedStatus: http.StatusBadRequest,
		},
		"error insert department from department service": {
			reqBody: rawMockDepartment,
			departmentService: map[string]testdata.FuncCall{
				"Create": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), &mockDepartment},
					Output: []interface{}{errors.New("unexpected error")},
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			for n, fn := range tc.departmentService {
				if fn.Called {
					mockDepartmentService.On(n, fn.Input...).Return(fn.Output...).Once()
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/departments", strings.NewReader(string(tc.reqBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			handler.AddDepartmentHandler(e, mockDepartmentService)

			e.ServeHTTP(rec, req)

			// var data interface{}
			// err := json.Unmarshal(rec.Body.Bytes(), &data)
			// require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestGet(t *testing.T) {
	e := testdata.GetEchoServer()
	e.Use(middleware.ErrorMiddleware())

	var mockDepartment domain.Department
	rawMockDepartment := testdata.GetGolden(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K")
	err := json.Unmarshal(rawMockDepartment, &mockDepartment)
	require.NoError(t, err)

	tests := map[string]struct {
		departmentID      string
		departmentService map[string]testdata.FuncCall
		expectedStatus    int
	}{
		"success": {
			departmentID: mockDepartment.ID,
			departmentService: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), mockDepartment.ID},
					Output: []interface{}{mockDepartment, nil},
				},
			},
			expectedStatus: http.StatusOK,
		},
		"not found": {
			departmentID: mockDepartment.ID,
			departmentService: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), mockDepartment.ID},
					Output: []interface{}{domain.Department{}, errors.New("not found")},
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		"error from department service": {
			departmentID: mockDepartment.ID,
			departmentService: map[string]testdata.FuncCall{
				"Get": testdata.FuncCall{
					Called: true,
					Input:  []interface{}{context.Background(), mockDepartment.ID},
					Output: []interface{}{domain.Department{}, errors.New("unexpected error")},
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	mockDepartmentService := new(mocks.DepartmentService)

	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			for name, fn := range testCase.departmentService {
				if fn.Called {
					mockDepartmentService.On(name, fn.Input...).Return(fn.Output...).Once()
				}
			}
			req := httptest.NewRequest(http.MethodGet, "/departments/"+testCase.departmentID, nil)

			rec := httptest.NewRecorder()
			handler.AddDepartmentHandler(e, mockDepartmentService)

			e.ServeHTTP(rec, req)

			res := rec.Result()

			require.Equal(t, testCase.expectedStatus, res.StatusCode)
		})
	}
}

func TestFetch(t *testing.T) {
	e := testdata.GetEchoServer()
	e.Use()

	mockDepartmentService := new(mocks.DepartmentService)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/departments", nil)

		rec := httptest.NewRecorder()
		handler.AddDepartmentHandler(e, mockDepartmentService)

		e.ServeHTTP(rec, req)

		res := rec.Result()

		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestUpdate(t *testing.T) {
	e := testdata.GetEchoServer()

	mockDepartmentService := new(mocks.DepartmentService)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/departments/123", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		handler.AddDepartmentHandler(e, mockDepartmentService)

		e.ServeHTTP(rec, req)

		res := rec.Result()

		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestDelete(t *testing.T) {
	e := testdata.GetEchoServer()

	mockDepartmentService := new(mocks.DepartmentService)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/departments/123", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		handler.AddDepartmentHandler(e, mockDepartmentService)

		e.ServeHTTP(rec, req)

		res := rec.Result()

		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})
}
