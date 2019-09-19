package mariadb

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	mgmysql "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/pkg/errors"
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
		dsnDB = "employee:employee-pass@tcp(localhost:3306)/employee?parseTime=1&loc=UTC&charset=utf8mb4&collation=utf8mb4_unicode_ci"
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

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("fail to run caller")
	}

	fmt.Println("++++++++ filename ++++++++")
	fmt.Printf("%+v\n", filename)
	fmt.Println("+++++++++++++++++")

	migrationPath := path.Join(path.Dir(filename), "migrations")

	fmt.Println("++++++++ migrations ++++++++")
	fmt.Printf("%+v\n", migrationPath)
	fmt.Println("+++++++++++++++++")

	m, err = migrate.NewWithDatabaseInstance("file://"+migrationPath, "mysql", driver)
	if err != nil {
		fmt.Println("++++++++ err migration++++++++")
		fmt.Printf("%+v\n", err)
		fmt.Println("+++++++++++++++++")
		return nil, err
	}

	err = m.Up()
	return
}

// TearDownSuite is a function to tear down test setup
func (d *DBSuite) TearDownSuite() {
	// require.NoError(d.T(), d.mg.Drop())
	require.NoError(d.T(), d.DB.Close())
}
