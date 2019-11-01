package mariadb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"

	repo "github.com/milhamhidayat/golang-clean-code-v2/department/repository/mariadb"
	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	mariadb "github.com/milhamhidayat/golang-clean-code-v2/driver/mariadb"
	ntime "github.com/milhamhidayat/golang-clean-code-v2/pkg/time"
	"github.com/milhamhidayat/golang-clean-code-v2/testdata"
)

type departmentSuite struct {
	mariadb.DBSuite
}

func TestDepartmentSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipped for short testing")
	}
	suite.Run(t, new(departmentSuite))
}

func (d *departmentSuite) SetupTest() {
	_, err := d.DB.Exec("TRUNCATE departments")
	require.NoError(d.T(), err)
}

func (d *departmentSuite) SeedDepartment(departments []domain.Department) (err error) {
	departmentRepo := repo.New(d.DB)
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
	departmentRepo := repo.New(d.DB)

	d.T().Run("success", func(t *testing.T) {
		var department domain.Department
		testdata.UnmarshallGoldenToJSON(t, "department-0ujssxh0cECutqzMgbtXSGnjorm", &department)

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
	departmentRepo := repo.New(d.DB)

	d.T().Run("success", func(t *testing.T) {
		var department domain.Department
		testdata.UnmarshallGoldenToJSON(t, "department-0ujssxh0cECutqzMgbtXSGnjorm", &department)

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
		_, err := departmentRepo.Get(context.Background(), "1")
		require.EqualError(t, err, domain.ErrNotFound.Error())
	})
}

func (d *departmentSuite) TestFetch() {
	departmentRepo := repo.New(d.DB)

	var department1, department2, department3, department4 domain.Department

	testdata.UnmarshallGoldenToJSON(d.T(), "department-0ujssxh0cECutqzMgbtXSGnjorm", &department1)
	testdata.UnmarshallGoldenToJSON(d.T(), "department-0ujsswThIGTUYm2K8FjOOfXtY1K", &department2)
	testdata.UnmarshallGoldenToJSON(d.T(), "department-0ujsszwN8NRY24YaXiTIE2VWDTS", &department3)
	testdata.UnmarshallGoldenToJSON(d.T(), "department-0ujsszgFvbiEr7CDgE3z8MAUPFt", &department4)

	departments := make([]domain.Department, 4)
	departments[0] = department1
	departments[1] = department2
	departments[2] = department3
	departments[3] = department4

	for i := range departments {
		date, err := ntime.GetLocalTime()
		require.NoError(d.T(), err)
		departments[i].CreatedTime = date
		departments[i].UpdatedTime = date
	}

	err := d.SeedDepartment(departments)
	require.NoError(d.T(), err)

	d.T().Run("success with ids", func(t *testing.T) {
		want := make([]domain.Department, 3)
		want[0] = departments[1]
		want[1] = departments[2]
		want[2] = departments[0]

		for i, v := range want {
			utcTime, err := ntime.ConvertToUTCTime(v.CreatedTime)
			require.NoError(t, err)

			want[i].CreatedTime = utcTime
			want[i].UpdatedTime = utcTime
		}

		depts, _, err := departmentRepo.Fetch(context.Background(), domain.DepartmentFilter{
			IDs: []string{want[0].ID, want[1].ID, want[2].ID},
		})

		require.Equal(t, want, depts)
		require.NoError(t, err)
	})

	d.T().Run("success with keyword", func(t *testing.T) {
		want := make([]domain.Department, 2)
		want[0] = departments[3]
		want[1] = departments[1]

		for i, v := range want {
			utcTime, err := ntime.ConvertToUTCTime(v.CreatedTime)
			require.NoError(t, err)

			want[i].CreatedTime = utcTime
			want[i].UpdatedTime = utcTime
		}

		expectedCursor := "MHVqc3N3VGhJR1RVWW0ySzhGak9PZlh0WTFL"

		depts, cur, err := departmentRepo.Fetch(context.Background(), domain.DepartmentFilter{
			Keyword: "Marketing",
		})

		require.Equal(t, want, depts)
		require.Equal(t, expectedCursor, cur)
		require.NoError(t, err)
	})

	d.T().Run("success with num", func(t *testing.T) {
		want := make([]domain.Department, 4)
		want[0] = departments[2]
		want[1] = departments[3]
		want[2] = departments[0]
		want[3] = departments[1]

		for i, v := range want {
			utcTime, err := ntime.ConvertToUTCTime(v.CreatedTime)
			require.NoError(t, err)

			want[i].CreatedTime = utcTime
			want[i].UpdatedTime = utcTime
		}

		expectedCursor := "MHVqc3N3VGhJR1RVWW0ySzhGak9PZlh0WTFL"

		depts, cur, err := departmentRepo.Fetch(context.Background(), domain.DepartmentFilter{
			Num: 4,
		})

		require.Equal(t, want, depts)
		require.Equal(t, expectedCursor, cur)
		require.NoError(t, err)
	})

	d.T().Run("success with num and cursor", func(t *testing.T) {
		var want []domain.Department

		expectedCursor := "MHVqc3N3VGhJR1RVWW0ySzhGak9PZlh0WTFL"
		depts, cur, err := departmentRepo.Fetch(context.Background(), domain.DepartmentFilter{
			Num:    4,
			Cursor: "MHVqc3N3VGhJR1RVWW0ySzhGak9PZlh0WTFL",
		})

		require.Equal(t, want, depts)
		require.Equal(t, expectedCursor, cur)
		require.NoError(t, err)
	})
}

func (d *departmentSuite) TestUpdate() {
	departmentRepo := repo.New(d.DB)

	var department domain.Department
	testdata.UnmarshallGoldenToJSON(d.T(), "department-0ujssxh0cECutqzMgbtXSGnjorm", &department)

	err := d.SeedDepartment([]domain.Department{department})
	require.NoError(d.T(), err)

	d.T().Run("success", func(t *testing.T) {
		newDepartment := domain.Department{
			ID:          department.ID,
			Name:        department.Name,
			Description: "this is description",
		}
		res, err := departmentRepo.Update(context.Background(), newDepartment)
		require.NoError(t, err)
		require.Equal(t, newDepartment.Description, res.Description)
	})
}

func (d *departmentSuite) TestDelete() {
	departmentRepo := repo.New(d.DB)

	var department1, department2 domain.Department
	testdata.UnmarshallGoldenToJSON(d.T(), "department-0ujssxh0cECutqzMgbtXSGnjorm", &department1)
	testdata.UnmarshallGoldenToJSON(d.T(), "department-0ujssxh0cECutqzMgbtXSGnjorm", &department2)

	d.T().Run("success", func(t *testing.T) {
		err := d.SeedDepartment([]domain.Department{department1})
		require.NoError(t, err)

		err = departmentRepo.Delete(context.TODO(), department1.ID)
		require.NoError(t, err)
	})
}
