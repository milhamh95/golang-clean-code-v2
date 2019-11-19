package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	"github.com/milhamhidayat/golang-clean-code-v2/pkg/validator"
)

type departmentHandler struct {
	service domain.DepartmentService
}

// AddDepartmentHandler adds the department handler
func AddDepartmentHandler(e *echo.Echo, service domain.DepartmentService) {
	if service == nil {
		panic("http: nil collection service")
	}

	handler := &departmentHandler{service}

	e.POST("/departments", handler.Insert)
	e.GET("/departments/:id", handler.Get)
	e.GET("/departments", handler.Fetch)
	e.PUT("/departments/:id", handler.Update)
	e.DELETE("/departments/:id", handler.Delete)
}

func (h departmentHandler) Insert(c echo.Context) error {
	ctx := c.Request().Context()

	var department domain.Department
	if err := c.Bind(&department); err != nil {
		return c.JSON(http.StatusBadRequest, "not ok")
	}

	if err := validator.Validate(department); err != nil {
		return c.JSON(http.StatusBadRequest, "not ok")
	}

	err := h.service.Create(ctx, &department)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, "ok")
}

func (h departmentHandler) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func (h departmentHandler) Fetch(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func (h departmentHandler) Update(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func (h departmentHandler) Delete(c echo.Context) error {
	return c.JSON(http.StatusNoContent, "ok")
}
