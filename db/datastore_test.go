package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit(t *testing.T) {
	var testDbConfig = DbConfig{
		DbSecrets: DbSecrets{
			RootPassword: "secret",
			MysqlUser:    "testuser",
			MysqlPwd:     "testpassword",
		},
		HostConfig: HostConfig{
			AutoRemove:    false,
			RestartPolicy: "no",
		},
		ExpireTime:   uint(240),
		DatabaseName: "csn_db",
	}
	connectionString, pool, resource := SetupDatbase(testDbConfig)

	port := resource.GetPort("3306/tcp")
	assert.Equal(t, "testuser:testpassword@(localhost:"+port+")/"+testDbConfig.DatabaseName+"?parseTime=true", *connectionString)
	defer Purge(pool, resource)
	db, initErr := InitDatabase(*connectionString)
	require.NoError(t, initErr, "Could not init the database")
	dbErr := db.Ping()
	require.NoError(t, dbErr, "Failed to ping database")
}
