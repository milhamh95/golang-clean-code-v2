package mariadb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	mariadb "github.com/milhamhidayat/golang-clean-code-v2/driver/mariadb"
	repo "github.com/milhamhidayat/golang-clean-code-v2/employee/repository/mariadb"
	ntime "github.com/milhamhidayat/golang-clean-code-v2/pkg/time"
	"github.com/milhamhidayat/golang-clean-code-v2/testdata"
)

type employeeSuite struct {
	mariadb.DBSuite
}

func TestEmployeeSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipped for short testing")
	}
	suite.Run(t, new(employeeSuite))
}

func (e *employeeSuite) SetupTest() {
	_, err := e.DB.Exec("TRUNCATE departments")
	require.NoError(e.T(), err)
	_, err = e.DB.Exec("TRUNCATE employees")
	require.NoError(e.T(), err)
}

func (e *employeeSuite) SeedEmployee(employees []domain.Employee) (err error) {
	employeeRepo := repo.New(e.DB)
	g, ctx := errgroup.WithContext(context.Background())

	for _, v := range employees {
		v := v
		g.Go(func() error {
			err := employeeRepo.Create(ctx, &v)
			return err
		})
	}

	if err = g.Wait(); err != nil {
		return
	}

	return
}

func (e *employeeSuite) SeedDepartment(departments []domain.Department) (err error) {
	c := context.Background()

	stmt, err := e.DB.PrepareContext(c, `INSERT INTO departments (id, name, description, created_time, updated_time) VALUES (?,?,?,?,?)`)
	require.NoError(e.T(), err)
	defer stmt.Close()

	g, ctx := errgroup.WithContext(c)
	for _, v := range departments {
		v := v
		g.Go(func() error {
			_, err := stmt.ExecContext(ctx, v.ID, v.Name, v.Description, v.CreatedTime, v.UpdatedTime)
			return err
		})
	}

	err = g.Wait()
	return
}

func (e *employeeSuite) TestCreate() {
	employeeRepo := repo.New(e.DB)

	e.T().Run("success", func(t *testing.T) {
		var employee domain.Employee
		testdata.UnmarshallGoldenToJSON(t, "employee-1S9XpJCvJbt1plvU36tAcJWS2ZW", &employee)

		err := employeeRepo.Create(context.Background(), &employee)
		require.NoError(t, err)

		newTime, err := ntime.ConvertToUTCTime(employee.CreatedTime)
		require.NoError(t, err)

		employee.CreatedTime = newTime
		employee.UpdatedTime = newTime

		emp, err := employeeRepo.Get(context.Background(), employee.ID)
		require.NoError(t, err)
		require.Equal(t, emp, employee)
	})
}

func (e *employeeSuite) TestGet() {
	employeeRepo := repo.New(e.DB)

	e.T().Run("success", func(t *testing.T) {
		var employee domain.Employee
		testdata.UnmarshallGoldenToJSON(t, "employee-1S9XpJCvJbt1plvU36tAcJWS2ZW", &employee)

		localTime, err := ntime.GetLocalTime()
		require.NoError(t, err)

		employee.CreatedTime = localTime
		employee.UpdatedTime = localTime

		err = e.SeedEmployee([]domain.Employee{employee})
		require.NoError(t, err)

		newTime, err := ntime.ConvertToUTCTime(employee.CreatedTime)
		require.NoError(t, err)

		employee.CreatedTime = newTime
		employee.UpdatedTime = newTime

		emp, err := employeeRepo.Get(context.Background(), employee.ID)
		require.NoError(t, err)
		require.Equal(t, emp, employee)
	})
}
