package db

import(
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Crud interface {
	Create()
	Read()
	Update()
	Delete()
}

func Init(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Could not create database %s!",err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Could not ping database %s!",err)
		return nil, err
	}
	return db, nil
}


