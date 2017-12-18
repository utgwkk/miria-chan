package main

import (
	"database/sql"
)

func NewMySQLConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	return db, err
}

func (m *MiriaClient) Sql() *sql.DB {
	if err := m.DB.Ping(); err != nil {
		db, err := NewMySQLConnection(m.DSN)
		if err != nil {
			panic(err)
		}
		m.DB = db
	}
	return m.DB
}
