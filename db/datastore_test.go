package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit(t *testing.T) {
	connectionString, pool, resource := SetuMySql("secret","testuser", "testpassword")
	port := resource.GetPort("3306")
	assert.Equal(t, "testuser:testpassword@(localhost:"+port+")/csn_db?parseTime=true", connectionString)
	defer Purge(pool, resource)
	db, initErr := Init(*connectionString)
	require.NoError(t, initErr, "Could not init the database")
	dbErr := db.Ping()
	require.NoError(t, dbErr, "Failed to ping database")
}
