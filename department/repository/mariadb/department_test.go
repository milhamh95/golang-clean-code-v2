package mariadb_test

import (
	"context"
	"fmt"
	"testing"

	repo "github.com/milhamhidayat/golang-clean-code-v2/department/repository/mariadb"
	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	mariadb "github.com/milhamhidayat/golang-clean-code-v2/driver/mariadb"
	"github.com/milhamhidayat/golang-clean-code-v2/testdata"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

// func (d *departmentSuite) SetupTest() {
// 	_, err := d.DB.Exec("TRUNCATE departments")
// 	require.NoError(d.T(), err)
// }

func (d *departmentSuite) TestCreate() {
	departmentRepo := repo.NewDepartmentRepository(d.DB)

	d.T().Run("success", func(t *testing.T) {
		var department domain.Department
		testdata.UnmarshallGoldenToJSON(t, "department-da7885d4-8504-47c7-9383-286a97faa14a", &department)

		err := departmentRepo.Create(context.Background(), department)
		require.NoError(t, err)

		dept, err := departmentRepo.Get(context.Background(), department.ID)
		require.NoError(t, err)

		fmt.Println("++++++++ res ++++++++")
		fmt.Printf("%+v\n", dept)
		fmt.Println("+++++++++++++++++")
	})
}
