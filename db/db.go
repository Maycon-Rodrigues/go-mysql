package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySql Driver
)

// Connection will provide a database connction
func Connection() (*sql.DB, error) {
	connectionString := "user:password@/db"
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Erro on open connection %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Erro to ping server %v", err)
		return nil, err
	}

	return db, nil
}
