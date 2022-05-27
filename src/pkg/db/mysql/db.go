package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var (
	Database *sql.DB
)

func Connect() {
	db, err := getDatabase()
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	Database = db
}

func getDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", viper.GetString("mysql_dsn"))
	return db, err
}
