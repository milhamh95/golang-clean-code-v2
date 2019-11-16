package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	handler "github.com/milhamhidayat/golang-clean-code-v2/department/delivery/http"
	"github.com/milhamhidayat/golang-clean-code-v2/domain/mocks"
	"github.com/milhamhidayat/golang-clean-code-v2/testdata"
)

func TestInsert(t *testing.T) {
	e := testdata.GetEchoServer()

	mockDepartmentService := new(mocks.DepartmentService)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/departments", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		handler.AddDepartmentHandler(e, mockDepartmentService)

		e.ServeHTTP(rec, req)

		res := rec.Result()

		require.Equal(t, http.StatusCreated, res.StatusCode)
	})
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
