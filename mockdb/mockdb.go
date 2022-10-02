package mockdb

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

var defaultMysqlConfig = &mysql.Config{
	User:                    "thaianh",
	Passwd:                  "thaianh1711",
	Net:                     "tcp",
	Addr:                    "localhost:3306",
	DBName:                  "mydb",
	AllowNativePasswords: true,
	MultiStatements: true,
}

var dbSkeleton *sql.DB

func OpenDB(config ...*mysql.Config) *sql.DB {
	if dbSkeleton != nil {
		return dbSkeleton
	}
	var cf *mysql.Config
	if len(config) > 0 {
		cf = config[0]
	} else {
		cf = defaultMysqlConfig
	}
	db, err := sql.Open("mysql", cf.FormatDSN())
	if err != nil {
		return nil
	}
	dbSkeleton = db
	return dbSkeleton
}