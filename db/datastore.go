package db

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"time"

	//_ "github.com/go-sql-driver/mysql"
	"log"
)

type Crud interface {
	Create(db *sql.DB)(sql.Result, error)
	Read(db *sql.DB)(error)
	Update(db *sql.DB)(sql.Result, error)
	Delete(db *sql.DB)(sql.Result, error)
}

var db *sql.DB

func SetuMySql(rootPassword string, mysqlUser string, mysqlPwd string) (*string, *dockertest.Pool, *dockertest.Resource) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "arngrimur/computersaysno_db",
		Tag:        "0.0.1",
		Env: []string{
			"MYSQL_ROOT_PASSWORD="+rootPassword,
			"MYSQL_USER=" + mysqlUser,
			"MYSQL_PASSWORD=" + mysqlPwd,
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = false
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	connectionsString := fmt.Sprintf(mysqlUser+":"+mysqlPwd+"@(localhost:%s)/csn_db?parseTime=true", resource.GetPort("3306/tcp"))
	log.Printf("Connections string: %s", connectionsString)
	resource.Expire(60) // Tell docker to hard kill the container

	pool.MaxWait = 60 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("mysql", connectionsString)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return &connectionsString, pool, resource
}

func Purge(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
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

