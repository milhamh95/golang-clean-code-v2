package mariadb_test

import (
	"context"
	"testing"

	repo "github.com/milhamhidayat/golang-clean-code-v2/department/repository/mariadb"
	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	mariadb "github.com/milhamhidayat/golang-clean-code-v2/driver/mariadb"
	ntime "github.com/milhamhidayat/golang-clean-code-v2/pkg/time"
	"github.com/milhamhidayat/golang-clean-code-v2/testdata"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"
)

type departmentSuite struct {
	mariadb.DBSuite
}

func TestDepartmentSuite(d *testing.T) {
	if testing.Short() {
		d.Skip("Skipped for short testing")
	}
	suite.Run(d, new(departmentSuite))
}

func (d *departmentSuite) SetupTest() {
	_, err := d.DB.Exec("TRUNCATE departments")
	require.NoError(d.T(), err)
}

func (d *departmentSuite) SeedDepartment(departments []domain.Department) (err error) {
	departmentRepo := repo.NewDepartmentRepository(d.DB)
	g, ctx := errgroup.WithContext(context.Background())

	for _, v := range departments {
		v := v
		g.Go(func() error {
			err := departmentRepo.Create(ctx, &v)
			return err
		})
	}

	if err = g.Wait(); err != nil {
		return
	}

	return
}

func (d *departmentSuite) TestCreate() {
	departmentRepo := repo.NewDepartmentRepository(d.DB)

	d.T().Run("success", func(t *testing.T) {
		var department domain.Department
		testdata.UnmarshallGoldenToJSON(t, "department-da7885d4-8504-47c7-9383-286a97faa14a", &department)

		date, err := ntime.GetLocalTime()
		require.NoError(t, err)
		department.CreatedTime = date
		department.UpdatedTime = date

		err = departmentRepo.Create(context.Background(), &department)
		require.NoError(t, err)

		utcTime, err := ntime.ConvertToUTCTime(date)
		require.NoError(t, err)

		department.CreatedTime = utcTime
		department.UpdatedTime = utcTime

		dept, err := departmentRepo.Get(context.Background(), department.ID)
		require.NoError(t, err)
		require.Equal(t, department, dept)
	})
}

func (d *departmentSuite) TestGet() {
	departmentRepo := repo.NewDepartmentRepository(d.DB)

	d.T().Run("success", func(t *testing.T) {
		var department domain.Department
		testdata.UnmarshallGoldenToJSON(t, "department-da7885d4-8504-47c7-9383-286a97faa14a", &department)

		date, err := ntime.GetLocalTime()
		require.NoError(t, err)
		department.CreatedTime = date
		department.UpdatedTime = date

		departments := []domain.Department{department}

		err = d.SeedDepartment([]domain.Department{department})
		require.NoError(t, err)

		utcTime, err := ntime.ConvertToUTCTime(date)
		require.NoError(t, err)

		departments[0].CreatedTime = utcTime
		departments[0].UpdatedTime = utcTime

		dept, err := departmentRepo.Get(context.Background(), department.ID)
		require.NoError(t, err)
		require.Equal(t, departments[0], dept)
	})

	d.T().Run("not found", func(t *testing.T) {
		expectedErr := errors.New("data is not found")
		_, err := departmentRepo.Get(context.Background(), "1")
		require.EqualError(t, err, expectedErr.Error())
	})
}
