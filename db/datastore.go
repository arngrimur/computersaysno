package db

import (
	"database/sql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"net"
	"net/url"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/lib/pq"
	"log"
)

type Crud interface {
	Create(db *sql.DB) (sql.Result, error)
	Read(db *sql.DB) error
	Update(db *sql.DB) (sql.Result, error)
	Delete(db *sql.DB) (sql.Result, error)
}

type DbSecrets struct {
	DatabaseUser     string
	DatabasePassword string
}
type HostConfig struct {
	AutoRemove    bool
	RestartPolicy string
}

type DbConfig struct {
	DbSecrets    DbSecrets
	HostConfig   HostConfig
	ExpireTime   uint
	DatabaseName string
}

var db *sql.DB
var connectionsString string

func SetupDatbase(dbConfig DbConfig) (*string, *dockertest.Pool, *dockertest.Resource) {

	pool, err := dockertest.NewPool("")
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	// pulls an image, creates a container based on it and runs it
	pool.MaxWait = time.Duration(dbConfig.ExpireTime) * time.Second
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	_, b, _, _ := runtime.Caller(0)
	workingDir := filepath.Dir(b)
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.1",
		Env: []string{
			"POSTGRES_USER=" + dbConfig.DbSecrets.DatabaseUser,
			"POSTGRES_PASSWORD=" + dbConfig.DbSecrets.DatabasePassword,
			"POSTGRES_DB=" + dbConfig.DatabaseName,
		},
		Mounts: []string{workingDir + "/mounts:/docker-entrypoint-initdb.d"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = dbConfig.HostConfig.AutoRemove
		config.RestartPolicy = docker.RestartPolicy{Name: dbConfig.HostConfig.RestartPolicy}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	resource.Expire(dbConfig.ExpireTime) // Tell docker to hard kill the container

	connectionsString = buildConnectionString(&dbConfig, resource)
	pool.MaxWait = time.Duration(dbConfig.ExpireTime) * time.Second
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	waitForConnection(err, pool, connectionsString)
	return &connectionsString, pool, resource
}

func buildConnectionString(dbConfig *DbConfig, resource *dockertest.Resource) string {
	pgUrl := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbConfig.DbSecrets.DatabaseUser, dbConfig.DbSecrets.DatabasePassword),
		Path:   dbConfig.DatabaseName,
	}
	q := pgUrl.Query()
	q.Add("sslmode", "disable")
	pgUrl.RawQuery = q.Encode()

	pgUrl.Host = resource.Container.NetworkSettings.IPAddress
	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		pgUrl.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}
	return pgUrl.String()
}
func waitForConnection(err error, pool *dockertest.Pool, connectionsString string) {
	if err = pool.Retry(wait); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}
func wait() error {
	db, err := sql.Open("postgres", connectionsString)
	if err != nil {
		return err
	}
	return db.Ping()
}

func Purge(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func InitDatabase(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
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
