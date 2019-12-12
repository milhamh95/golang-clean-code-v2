package http

import (
	"github.com/friendsofgo/errors"
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
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := validator.Validate(department); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := h.service.Create(ctx, &department)
	if err != nil {
		return errors.Wrap(err, "failed to insert a department")
	}

	return c.JSON(http.StatusCreated, department)
}

func (h departmentHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	departmentID := c.Param("id")

	department, err := h.service.Get(ctx, departmentID)
	if err != nil {
		return errors.Wrap(err, "failed get a department")
	}

	return c.JSON(http.StatusOK, department)
}

func (h departmentHandler) Fetch(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func (h departmentHandler) Update(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func (h departmentHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	departmentID := c.Param("id")

	err := h.service.Delete(ctx, departmentID)
	if err != nil {
		return errors.Wrap(err, "failed delete a department")
	}
	return c.NoContent(http.StatusNoContent)
}
