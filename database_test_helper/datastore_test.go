package database_test_helper

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit(t *testing.T) {

	connectionString, pool, resource := SetupDatbase()

	t.Cleanup(func() {
		Purge(pool, resource)
	})
	assert.Regexp(t, "^postgres://testuser:testpassword@((\\d){1,3}\\.){3}\\d{1,3}:?(\\d){0,5}\\/"+TestDbConfig.DatabaseName+"\\?sslmode=disable$", *connectionString)
	db, initErr := InitDatabase(*connectionString)
	require.NoError(t, initErr, "Could not init the database")
	dbErr := db.Ping()
	require.NoError(t, dbErr, "Failed to ping database")
}
