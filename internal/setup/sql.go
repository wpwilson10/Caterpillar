package setup

import (
	"os"

	"github.com/jmoiron/sqlx"
	// blank import needed for sqlx to work
	_ "github.com/lib/pq"
)

// SQL initalizes a database connection based off .env file connection parameters
func SQL() *sqlx.DB {
	var host string = "host=" + os.Getenv("SQL_HOST")
	var port string = "port=" + os.Getenv("SQL_PORT")
	var user string = "user=" + os.Getenv("SQL_USER")
	var password string = "password=" + os.Getenv("SQL_PASSWORD")
	var dbname string = "dbname=" + os.Getenv("SQL_DB")
	var connectionString = host + " " + port + " " + user + " " + password + " " + dbname + " " + "sslmode=disable"

	// this Pings the database trying to connect, panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		LogCommon(err).Fatal("Connecting to SQL database")
	}

	return db
}
