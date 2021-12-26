package db

import (
	"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	"log"
)

type Crud interface {
	Create(db *sql.DB)(sql.Result, error)
	Read(db *sql.DB)(error)
	Update(db *sql.DB)(sql.Result, error)
	Delete(db *sql.DB)(sql.Result, error)
}

func Init(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connectionString)
	db.SetMaxIdleConns(0)
	if err != nil {
		log.Fatal("Could not create database!", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Could not ping database!", err)
		return nil, err
	}
	return db, nil
}

