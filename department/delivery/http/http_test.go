package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	handler "github.com/milhamhidayat/golang-clean-code-v2/department/delivery/http"
	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	"github.com/milhamhidayat/golang-clean-code-v2/domain/mocks"
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/departments", strings.NewReader(string(tc.reqBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			handler.AddDepartmentHandler(e, mockDepartmentService)

			e.ServeHTTP(rec, req)

			var data interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &data)
			require.NoError(t, err)

			fmt.Println("========  ========")
			fmt.Printf("%+v\n", data)
			fmt.Println("=================")

			require.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestGet(t *testing.T) {
	e := testdata.GetEchoServer()

	mockDepartmentService := new(mocks.DepartmentService)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/departments/123", nil)

		rec := httptest.NewRecorder()
		handler.AddDepartmentHandler(e, mockDepartmentService)

		e.ServeHTTP(rec, req)

		res := rec.Result()

		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestFetch(t *testing.T) {
	e := testdata.GetEchoServer()

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
