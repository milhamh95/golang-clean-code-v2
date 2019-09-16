package mariadb

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	mgmysql "github.com/golang-migrate/migrate/database/mysql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// DBSuite is a test suite for maria db
type DBSuite struct {
	suite.Suite
	DB *sql.DB
	mg *migrate.Migrate
}

// SetupSuite is a function to setup test suite for maria db
func (d *DBSuite) SetupSuite() {
	dsnDB := os.Getenv("MYSQL_TEST")
	if dsnDB == "" {
		dsnDB = "employee:employee-pass@tcp(mysql:3306)/employee?parseTime=1&loc=UTC&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	}

	db, err := sql.Open("mysql", dsnDB)
	require.NoError(d.T(), err)
	require.NotNil(d.T(), db)

	d.mg, err = MigrateDB(db)
	require.NoError(d.T(), err)
	d.DB = db

}

// MigrateDB is a function to migrate a db
func MigrateDB(db *sql.DB) (m *migrate.Migrate, err error) {
	driver, err := mgmysql.WithInstance(db, &mgmysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err = migrate.NewWithDatabaseInstance("file://migrations/", "mysql", driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	return
}

// TearDownSuite is a function to tear down test setup
func (d *DBSuite) TearDownSuite() {
	require.NoError(d.T(), d.mg.Drop())
	require.NoError(d.T(), d.DB.Close())
}
