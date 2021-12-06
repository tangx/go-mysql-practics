package dbconn

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sirupsen/logrus"
)

func InitDB() *sql.DB {
	dsn := `root:Mysql12345@tcp(127.0.0.1:3306)/mysqlpractics`
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
