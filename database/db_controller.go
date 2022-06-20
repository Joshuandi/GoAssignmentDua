package database

import "database/sql"

const (
	Host     = "localhost"
	Port     = 5432
	User     = "postgres"
	Password = "admin"
	Dbname   = "Assignment2"
)

var (
	Db  *sql.DB
	Err error
)
