package mariadb_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
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
	_, err := e.DB.Exec("DELETE FROM employees")
	require.NoError(e.T(), err)
	_, err = e.DB.Exec("DELETE FROM departments")
	require.NoError(e.T(), err)
}

func (e *employeeSuite) SeedEmployee(employees []domain.Employee) (err error) {
	e.T().Helper()
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
	e.T().Helper()
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
		var (
			department domain.Department
			employee   domain.Employee
		)
		testdata.UnmarshallGoldenToJSON(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &department)
		testdata.UnmarshallGoldenToJSON(t, "employee-1S9XpJCvJbt1plvU36tAcJWS2ZW", &employee)

		err := e.SeedDepartment([]domain.Department{department})
		require.NoError(t, err)

		err = employeeRepo.Create(context.Background(), &employee)
		require.NoError(t, err)

		newTime, err := ntime.ConvertToUTCTime(employee.CreatedTime)
		require.NoError(t, err)

		employee.CreatedTime = newTime
		employee.UpdatedTime = newTime

		emp, err := employeeRepo.Get(context.Background(), employee.ID)
		require.NoError(t, err)
		require.Equal(t, employee, emp)
	})
}

func (e *employeeSuite) TestGet() {
	employeeRepo := repo.New(e.DB)

	e.T().Run("success", func(t *testing.T) {
		var (
			department domain.Department
			employee   domain.Employee
		)
		testdata.UnmarshallGoldenToJSON(t, "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &department)
		testdata.UnmarshallGoldenToJSON(t, "employee-1S9XpJCvJbt1plvU36tAcJWS2ZW", &employee)

		err := e.SeedDepartment([]domain.Department{department})
		require.NoError(t, err)

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
		require.Equal(t, employee, emp)
	})

	e.T().Run("not found", func(t *testing.T) {
		expectedErr := errors.New("employee is not found: 1")
		_, err := employeeRepo.Get(context.Background(), "1")
		require.EqualError(t, err, expectedErr.Error())
	})
}

func (e *employeeSuite) TestFetch() {
	employeeRepo := repo.New(e.DB)

	var employee1, employee2 domain.Employee

	testdata.UnmarshallGoldenToJSON(e.T(), "employee-1S9XpJCvJbt1plvU36tAcJWS2ZW", &employee1)
	testdata.UnmarshallGoldenToJSON(e.T(), "employee-1SYxHnSCbFCxLr7zUxk5j8cB0Cr", &employee2)

	employees := make([]domain.Employee, 2)
	employees[0] = employee1
	employees[1] = employee2

	for i := range employees {
		date, err := ntime.GetLocalTime()
		require.NoError(e.T(), err)
		employees[i].CreatedTime = date
		employees[i].UpdatedTime = date
	}

	err := e.SeedEmployee(employees)
	require.NoError(e.T(), err)

	e.T().Run("success with ids", func(t *testing.T) {
		expectedEmployees := make([]domain.Employee, 2)
		expectedEmployees[0] = employees[1]
		expectedEmployees[1] = employees[0]

		for i, v := range expectedEmployees {
			utcTime, err := ntime.ConvertToUTCTime(v.CreatedTime)
			require.NoError(t, err)

			expectedEmployees[i].CreatedTime = utcTime
			expectedEmployees[i].UpdatedTime = utcTime
		}

		emps, cursor, err := employeeRepo.Fetch(context.Background(), domain.EmployeeFilter{
			IDs: []string{expectedEmployees[0].ID, expectedEmployees[1].ID},
		})

		require.NoError(t, err)
		require.Equal(t, expectedEmployees, emps)
		require.Equal(t, "", cursor)
	})

	e.T().Run("success with dept ids", func(t *testing.T) {
		expectedEmployees := make([]domain.Employee, 1)
		expectedEmployees[0] = employees[1]

		for i, v := range expectedEmployees {
			utcTime, err := ntime.ConvertToUTCTime(v.CreatedTime)
			require.NoError(t, err)

			expectedEmployees[i].CreatedTime = utcTime
			expectedEmployees[i].UpdatedTime = utcTime

			emps, cursor, err := employeeRepo.Fetch(context.Background(), domain.EmployeeFilter{
				IDs: []string{expectedEmployees[0].ID},
			})

			require.NoError(t, err)
			require.Equal(t, expectedEmployees, emps)
			require.Equal(t, "", cursor)
		}
	})

	e.T().Run("success with num", func(t *testing.T) {
		expectedEmployees := make([]domain.Employee, 2)
		expectedEmployees[0] = employees[1]
		expectedEmployees[1] = employees[0]

		for i, v := range expectedEmployees {
			utcTime, err := ntime.ConvertToUTCTime(v.CreatedTime)
			require.NoError(t, err)

			expectedEmployees[i].CreatedTime = utcTime
			expectedEmployees[i].UpdatedTime = utcTime
		}

		emps, cursor, err := employeeRepo.Fetch(context.Background(), domain.EmployeeFilter{
			Num: 2,
		})

		require.NoError(t, err)
		require.Equal(t, expectedEmployees, emps)
		require.Equal(t, "MVM5WHBKQ3ZKYnQxcGx2VTM2dEFjSldTMlpX", cursor)
	})

	e.T().Run("success with second page using num and cursor", func(t *testing.T) {
		emps, cursor, err := employeeRepo.Fetch(context.Background(), domain.EmployeeFilter{
			Num:    2,
			Cursor: "MVM5WHBKQ3ZKYnQxcGx2VTM2dEFjSldTMlpX",
		})

		require.NoError(t, err)
		require.Equal(t, []domain.Employee{}, emps)
		require.Equal(t, "MVM5WHBKQ3ZKYnQxcGx2VTM2dEFjSldTMlpX", cursor)
	})

	e.T().Run("success with keyword", func(t *testing.T) {
		expectedEmployees := make([]domain.Employee, 1)
		expectedEmployees[0] = employees[1]

		for i, v := range expectedEmployees {
			utcTime, err := ntime.ConvertToUTCTime(v.CreatedTime)
			require.NoError(t, err)

			expectedEmployees[i].CreatedTime = utcTime
			expectedEmployees[i].UpdatedTime = utcTime
		}

		emps, cursor, err := employeeRepo.Fetch(context.Background(), domain.EmployeeFilter{
			Keyword: "casey",
		})

		require.NoError(t, err)
		require.Equal(t, expectedEmployees, emps)
		require.Equal(t, "MVNZeEhuU0NiRkN4THI3elV4azVqOGNCMENy", cursor)
	})
}
