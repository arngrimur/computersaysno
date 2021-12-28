package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit(t *testing.T) {
	var testDbConfig = DbConfig{
		DbSecrets: DbSecrets{
			DatabaseUser:     "testuser",
			DatabasePassword: "testpassword",
		},
		HostConfig: HostConfig{
			AutoRemove:    true,
			RestartPolicy: "no",
		},
		ExpireTime:   uint(120),
		DatabaseName: "csn_db",
	}
	connectionString, pool, resource := SetupDatbase(testDbConfig)

	port := resource.GetPort("5432/tcp")
	t.Cleanup(func() {
		Purge(pool, resource)
	})
	assert.Regexp(t, "postgres://testuser:testpassword@(\\d{1-3}\\.){4}:"+port+"/"+testDbConfig.DatabaseName+"?sslmode=disable", *connectionString)
	db, initErr := InitDatabase(*connectionString)
	require.NoError(t, initErr, "Could not init the database")
	dbErr := db.Ping()
	require.NoError(t, dbErr, "Failed to ping database")
}
