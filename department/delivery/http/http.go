package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/friendsofgo/errors"

	"github.com/labstack/echo/v4"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	"github.com/milhamhidayat/golang-clean-code-v2/pkg/md5"
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
	ctx := c.Request().Context()

	keyword := c.QueryParam("keyword")
	cursor := c.QueryParam("cursor")

	ids := make([]string, 0)
	paramIDs := c.QueryParam("ids")
	if paramIDs != "" {
		ids = strings.Split(paramIDs, ",")
	}

	num := 20
	if numStr := c.QueryParam("num"); numStr != "" {
		var err error
		if num, err = strconv.Atoi(numStr); err != nil {
			err = fmt.Errorf("num query-param is not valid. Got error when parsing value: %v", err)
			return domain.ConstraintErrorf("%s", err)
		}
	}

	res, nextCursor, err := h.service.Fetch(ctx, domain.DepartmentFilter{IDs: ids, Keyword: keyword, Num: num, Cursor: cursor})
	if err != nil {
		return errors.Wrap(err, "error fetch departments")
	}

	if len(res) > 0 {
		eTag := ""
		if eTag, err = md5.Generate(res[0].ID); err != nil {
			return errors.Wrap(err, "error generate departments eTag")
		}

		ifNoneMatch := c.Request().Header.Get("If-None-Match")
		if eTag != "" && ifNoneMatch != "" && strings.Contains(ifNoneMatch, eTag) {
			return c.NoContent(http.StatusNotModified)
		}

		c.Response().Header().Set("ETag", "W/"+eTag)
		c.Response().Header().Set("X-Cursor", nextCursor)
	}

	return c.JSON(http.StatusOK, res)
}

func (h departmentHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	departmentID := c.Param("id")

	var department domain.Department
	if err := c.Bind(&department); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := validator.Validate(department); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	department.ID = departmentID

	res, err := h.service.Update(ctx, department)
	if err != nil {
		return errors.Wrap(err, "failed to delete a department")
	}

	return c.JSON(http.StatusOK, res)
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
