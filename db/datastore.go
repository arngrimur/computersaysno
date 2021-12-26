package db

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Crud interface {
	Create(db *sql.DB) (sql.Result, error)
	Read(db *sql.DB) error
	Update(db *sql.DB) (sql.Result, error)
	Delete(db *sql.DB) (sql.Result, error)
}

type DbSecrets struct {
	RootPassword string
	MysqlUser    string
	MysqlPwd     string
}
type HostConfig struct {
	AutoRemove    bool
	RestartPolicy string
}

type DbConfig struct {
	DbSecrets  DbSecrets
	HostConfig HostConfig
	ExpireTime uint
}

var db *sql.DB

func SetuMySql(dbConfig DbConfig) (*string, *dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	// pulls an image, creates a container based on it and runs it
	pool.MaxWait = time.Duration(dbConfig.ExpireTime) * time.Second
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "arngrimur/computersaysno_db",
		Tag:        "0.0.1",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=" + dbConfig.DbSecrets.RootPassword,
			"MYSQL_USER=" + dbConfig.DbSecrets.MysqlUser,
			"MYSQL_PASSWORD=" + dbConfig.DbSecrets.MysqlPwd,
		},
		// set AutoRemove to true so that stopped container goes away by itself
	}, func(config *docker.HostConfig) {
		config.AutoRemove = dbConfig.HostConfig.AutoRemove
		config.RestartPolicy = docker.RestartPolicy{Name: dbConfig.HostConfig.RestartPolicy}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	resource.Expire(dbConfig.ExpireTime) // Tell docker to hard kill the container
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	connectionsString := fmt.Sprintf(dbConfig.DbSecrets.MysqlUser+":"+dbConfig.DbSecrets.MysqlPwd+"@(localhost:%s)/csn_db?parseTime=true", resource.GetPort("3306/tcp"))
	log.Printf("Connections string: %s", connectionsString)

	if err = pool.Retry(func() error {
		db, err = sql.Open("mysql", connectionsString)
		if err != nil {
			return err
		}
		time.Sleep(30 * time.Second)
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
