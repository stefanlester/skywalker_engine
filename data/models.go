package data

import (
	db2 "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"

	"database/sql"
	"os"
)

var db *sql.DB
var upper db2.Session

// OpenDB opens a database connection and returns a pointer to sql.DB

type Models struct {
	// any models inserted here (and in new function)
	// are easily accessible throughtout the entire application
}

func New(databasePool *sql.DB) Models {
	db = databasePool

	if os.Getenv("DATABASE_TYPE") == "mysql" || os.Getenv("DATABASE_TYPE") == "mariadb" {
		upper, _ = mysql.New(databasePool)
	} else {
		upper, _ = postgresql.New(databasePool)
	}

	return Models{}
}
