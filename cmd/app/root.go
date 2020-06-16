package main

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	deptRepo "github.com/milhamhidayat/golang-clean-code-v2/department/repository/mariadb"
	deptService "github.com/milhamhidayat/golang-clean-code-v2/department/service"
	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	"github.com/milhamhidayat/golang-clean-code-v2/pkg/env"
)

var (
	departmentRepository domain.DepartmentRepository
	departmentService    domain.DepartmentService
)

var rootCmd = &cobra.Command{
	Use:   "employee",
	Short: "Employee is an API for managing employees",
}

func init() {
	cobra.OnInitialize(initApp)
}

// Execute the main function
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func initApp() {
	/**
	 * MYSQL Conf
	 */
	dsnMysql := env.Get("MYSQL_URI")
	db, err := sql.Open("mysql", dsnMysql)
	if err != nil {
		log.Fatal("can't open mysql connection to: %s, got err: %v", dsnMysql, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("can't connect to mysql db, err: %v", err)
	}

	mysqlMaxIdleCon, err := strconv.Atoi(env.Get("MYSQL_MAX_IDLE_CONNECTION"))
	if err != nil {
		log.Fatal("MYSQL_MAX_IDLE_CONNECTION is not well-set")
	}
	db.SetMaxIdleConns(mysqlMaxIdleCon)

	mysqlMaxOpenCon, err := strconv.Atoi(env.Get("MYSQL_MAX_OPEN_CONNECTION"))
	if err != nil {
		log.Fatal("MYSQL_MAX_OPEN_CONNECTION is not well-set")
	}
	db.SetMaxOpenConns(mysqlMaxOpenCon)

	mysqlMaxConnLifetime, err := strconv.Atoi(env.Get("MYSQL_CONNECTION_LIFETIME_M"))
	if err != nil {
		log.Fatal("MYSQL_CONNECTION_LIFETIME_M is not well-set")
	}
	db.SetConnMaxLifetime(time.Minute * time.Duration(mysqlMaxConnLifetime))

	/**
	 * Context Timeout
	 */
	// t, err := strconv.ParseInt(env.Get("CONTEXT_TIMEOUT_MS"), 10, 16)
	// if err != nil {
	// 	log.Fatal("CONTEXT_TIMEOUT_MS is not well-set")
	// }
	// contextTimeout := time.Duration(t) * time.Millisecond

	/**
	 * Department
	 */
	departmentRepository = deptRepo.New(db)
	departmentService = deptService.New(departmentRepository)
}
