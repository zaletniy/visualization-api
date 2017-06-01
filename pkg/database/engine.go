package db

import (
	"fmt"
	// import mysql driver for side-effect required for xorm package
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	engine *xorm.Engine
)

// InitializeEngine initializes connection to db
func InitializeEngine(mysqlUsername, mysqlPassword, mysqlHost,
	mysqlDatabaseName string, mysqlPort int) error {
	var err error
	// connection parameter has to look like
	// user:password@tcp(host:port)/dbname))
	engine, err = xorm.NewEngine("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		mysqlUsername,
		mysqlPassword,
		mysqlHost,
		mysqlPort,
		mysqlDatabaseName,
	))
	return err
}

// GetEngine returns database connection
func GetEngine() *xorm.Engine {
	return engine
}
