package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
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
	testdata.UnmarshallGoldenToJSON(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &mockDepartment)
	rawMockDepartment := testdata.GetGolden(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K")

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
			mockDepartmentService := new(mocks.DepartmentService)

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

			mockDepartmentService.AssertExpectations(t)

			require.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestGet(t *testing.T) {
	e := testdata.GetEchoServer()
	e.Use(middleware.ErrorMiddleware())

	var mockDepartment domain.Department
	testdata.UnmarshallGoldenToJSON(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &mockDepartment)

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
					Output: []interface{}{domain.Department{}, domain.ErrNotFound},
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

	for testName, testCase := range tests {
		mockDepartmentService := new(mocks.DepartmentService)
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

			mockDepartmentService.AssertExpectations(t)

			res := rec.Result()

			require.Equal(t, testCase.expectedStatus, res.StatusCode)
		})
	}
}

func TestFetch(t *testing.T) {
	var departments []domain.Department
	testdata.UnmarshallGoldenToJSON(t, "departments", &departments)

	var engineerDepartment domain.Department
	testdata.UnmarshallGoldenToJSON(t, "department-0ujssxh0cECutqzMgbtXSGnjorm", &engineerDepartment)
	engineerDepartments := []domain.Department{engineerDepartment}

	e := testdata.GetEchoServer()
	e.Use(middleware.ErrorMiddleware())

	tests := map[string]struct {
		departmentService  testdata.FuncCall
		target             string
		ifNoneMatch        string
		expectedStatusCode int
		expectedCursor     string
		expectedETag       string
	}{
		"success with num": {
			departmentService: testdata.FuncCall{
				Called: true,
				Input: []interface{}{mock.Anything, domain.DepartmentFilter{
					IDs:     []string{},
					Keyword: "",
					Num:     20,
					Cursor:  "",
				}},
				Output: []interface{}{departments, "next-cursor", nil},
			},
			target:             "/departments",
			expectedStatusCode: http.StatusOK,
			expectedCursor:     "next-cursor",
			expectedETag:       "W/d60c95250ff1839e44dc74409f2b6c63",
		},
		"success with keyword": {
			departmentService: testdata.FuncCall{
				Called: true,
				Input: []interface{}{mock.Anything, domain.DepartmentFilter{
					IDs:     []string{},
					Keyword: "engineer",
					Num:     20,
					Cursor:  "",
				}},
				Output: []interface{}{engineerDepartments, "next-cursor", nil},
			},
			target:             "/departments?keyword=engineer",
			expectedStatusCode: http.StatusOK,
			expectedCursor:     "next-cursor",
			expectedETag:       "W/cbd902cb9cd45600989fdca27dbdbbe0",
		},
		"success with ids": {
			departmentService: testdata.FuncCall{
				Called: true,
				Input: []interface{}{mock.Anything, domain.DepartmentFilter{
					IDs:     []string{"0ujssxh0cECutqzMgbtXSGnjorm"},
					Keyword: "",
					Num:     20,
					Cursor:  "",
				}},
				Output: []interface{}{engineerDepartments, "", nil},
			},
			target:             "/departments?ids=0ujssxh0cECutqzMgbtXSGnjorm",
			expectedStatusCode: http.StatusOK,
			expectedCursor:     "",
			expectedETag:       "W/cbd902cb9cd45600989fdca27dbdbbe0",
		},
		"success with etag": {
			departmentService: testdata.FuncCall{
				Called: true,
				Input: []interface{}{mock.Anything, domain.DepartmentFilter{
					IDs:     []string{},
					Keyword: "",
					Num:     20,
					Cursor:  "",
				}},
				Output: []interface{}{departments, "next-cursor", nil},
			},
			target:             "/departments",
			ifNoneMatch:        "W/d60c95250ff1839e44dc74409f2b6c63",
			expectedStatusCode: http.StatusNotModified,
			expectedCursor:     "",
			expectedETag:       "",
		},
		"with bad param": {
			departmentService: testdata.FuncCall{
				Called: false,
			},
			target:             "/departments?num=xxxx",
			expectedStatusCode: http.StatusBadRequest,
		},
		"with unexpected error": {
			departmentService: testdata.FuncCall{
				Called: true,
				Input: []interface{}{mock.Anything, domain.DepartmentFilter{
					IDs:     []string{},
					Keyword: "",
					Num:     20,
					Cursor:  "",
				}},
				Output: []interface{}{[]domain.Department{}, "", errors.New("unexpected error")},
			},
			target:             "/departments",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			departmentServiceMock := new(mocks.DepartmentService)
			if test.departmentService.Called {
				departmentServiceMock.On("Fetch", test.departmentService.Input...).
					Return(test.departmentService.Output...).Once()
			}

			req := httptest.NewRequest(http.MethodGet, test.target, nil)
			if test.ifNoneMatch != "" {
				req.Header.Set("If-None-Match", test.ifNoneMatch)
			}
			rec := httptest.NewRecorder()

			handler.AddDepartmentHandler(e, departmentServiceMock)
			e.ServeHTTP(rec, req)

			departmentServiceMock.AssertExpectations(t)

			require.Equal(t, test.expectedCursor, rec.Header().Get("X-Cursor"))
			require.Equal(t, test.expectedETag, rec.Header().Get("ETag"))
			require.Equal(t, test.expectedStatusCode, rec.Code)
		})
	}
}

func TestUpdate(t *testing.T) {
	e := testdata.GetEchoServer()
	e.Use(middleware.ErrorMiddleware())

	var department domain.Department
	testdata.UnmarshallGoldenToJSON(t, "department-0ujssxh0cECutqzMgbtXSGnjorm", &department)

	deptReq := department
	deptReq.CreatedTime = time.Time{}
	deptReq.UpdatedTime = time.Time{}
	deptReqJSON, err := json.Marshal(deptReq)
	require.NoError(t, err)

	tests := map[string]struct {
		reqBody          []byte
		departmentID     string
		depatmentService testdata.FuncCall
		expectedStatus   int
	}{
		"success": {
			reqBody:      deptReqJSON,
			departmentID: "0ujssxh0cECutqzMgbtXSGnjorm",
			depatmentService: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, deptReq},
				Output: []interface{}{department, nil},
			},
			expectedStatus: http.StatusOK,
		},
		"not found": {
			reqBody:      deptReqJSON,
			departmentID: "0ujssxh0cECutqzMgbtXSGnjorm",
			depatmentService: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, deptReq},
				Output: []interface{}{domain.Department{}, domain.ErrNotFound},
			},
			expectedStatus: http.StatusNotFound,
		},
		"unexpected error": {
			reqBody:      deptReqJSON,
			departmentID: "0ujssxh0cECutqzMgbtXSGnjorm",
			depatmentService: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, deptReq},
				Output: []interface{}{domain.Department{}, errors.New("unexpected error")},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			departmentServiceMock := new(mocks.DepartmentService)
			if test.depatmentService.Called {
				departmentServiceMock.On("Update", test.depatmentService.Input...).
					Return(test.depatmentService.Output...).Once()
			}

			handler.AddDepartmentHandler(e, departmentServiceMock)

			req := httptest.NewRequest(http.MethodPut, "/departments/"+test.departmentID, strings.NewReader(string(test.reqBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			departmentServiceMock.AssertExpectations(t)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestDelete(t *testing.T) {
	e := testdata.GetEchoServer()
	e.Use(middleware.ErrorMiddleware())

	tests := map[string]struct {
		departmentID      string
		departmentService testdata.FuncCall
		expectedStatus    int
	}{
		"success": {
			departmentID: "0ujssxh0cECutqzMgbtXSGnjorm",
			departmentService: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, "0ujssxh0cECutqzMgbtXSGnjorm"},
				Output: []interface{}{nil},
			},
			expectedStatus: http.StatusNoContent,
		},
		"not found": {
			departmentID: "0ujssxh0cECutqzMgbtXSGnjorm",
			departmentService: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, "0ujssxh0cECutqzMgbtXSGnjorm"},
				Output: []interface{}{domain.ErrNotFound},
			},
			expectedStatus: http.StatusNotFound,
		},
		"unexpected error": {
			departmentID: "0ujssxh0cECutqzMgbtXSGnjorm",
			departmentService: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, "0ujssxh0cECutqzMgbtXSGnjorm"},
				Output: []interface{}{errors.New("unexpected error")},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			mockDepartmentService := new(mocks.DepartmentService)
			if test.departmentService.Called {
				mockDepartmentService.On("Delete", test.departmentService.Input...).
					Return(test.departmentService.Output...).Once()
			}

			req := httptest.NewRequest(http.MethodDelete, "/departments/"+test.departmentID, nil)
			rec := httptest.NewRecorder()
			handler.AddDepartmentHandler(e, mockDepartmentService)

			e.ServeHTTP(rec, req)

			mockDepartmentService.AssertExpectations(t)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}
