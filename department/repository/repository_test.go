package repository_test

import (
	"testing"

	"github.com/milhamhidayat/golang-clean-code-v2/driver/mariadb"
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

func (d *departmentSuite) SetupTest() {
	_, err := d.DB.Exec("TRUNCATE departments")
	require.NoError(d.T(), err)
}
