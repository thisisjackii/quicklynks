# File: backend/db/db.go
package db

import (
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"          // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// InitDB initializes and returns a database connection pool.
func InitDB(dataSourceName string) (*sqlx.DB, error) {
	var driverName string
	if strings.HasPrefix(dataSourceName, "postgres://") {
		driverName = "postgres"
	} else {
		driverName = "sqlite3"
	}

	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// ApplySchema reads and executes the schema.sql file.
func ApplySchema(db *sqlx.DB) error {
	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(schema))
	return err
}