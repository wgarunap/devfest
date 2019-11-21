package adapters

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// NewDbConnection initialize and returns a new database connection
func NewDbConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:123@/phonebook")

	if err != nil {
		panic(err)
	}

	return db
}
