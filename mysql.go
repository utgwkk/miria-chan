package main

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(hostname, databaseName, username, password string) (*sql.DB, error) {
	config := &mysql.Config{
		User:   username,
		Passwd: password,
		DBName: databaseName,
		Addr:   hostname,
	}
	dsn := config.FormatDSN()
	db, err := sql.Open("mysql", dsn)
	return db, err
}
