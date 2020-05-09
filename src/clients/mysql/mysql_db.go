package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dzikrisyafi/kursusvirtual_utils-go/logger"
	"github.com/go-sql-driver/mysql"
)

var (
	client *sql.DB

	username = os.Getenv("mysql_username")
	password = os.Getenv("mysql_password")
	host     = os.Getenv("mysql_host")
	schema   = os.Getenv("mysql_schema1")
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema,
	)

	var err error
	client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err = client.Ping(); err != nil {
		panic(err)
	}

	mysql.SetLogger(logger.GetLogger())
	log.Println("database successfully configured")
}

func DbConn() *sql.DB {
	return client
}
