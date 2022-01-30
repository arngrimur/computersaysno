package db

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

func InitDatabase() (*sql.DB, error) {
	connString := createConnectionString()
	db, openErr := sql.Open("postgres", connString)
	if openErr != nil {
		log.Fatalf("Can create connection with database at '%s'", connString)
		return nil, openErr
	}
	driver, driverErr := postgres.WithInstance(db, &postgres.Config{})
	if driverErr != nil {
		log.Fatalf("Can not instanciate database")
		return nil, driverErr
	}
	m, migrateErr := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"csn_db", driver)
	if migrateErr != nil {
		log.Fatalf("Can not migrate datbase.")
		return nil, migrateErr
	}
	upErr := m.Up()
	if upErr != nil {
		return nil, upErr
	} // or m.Step(2) if you want to explicitly set the number of migrations to run
	return db, nil
}

func createConnectionString() string {
	log.Fatalf("Implement me!!! Grab from kubenetes secets")
	return "jdbc:postgresql://localhost:55144/csn_db?sslmode=disable"
}
