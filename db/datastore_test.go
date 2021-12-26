package db

import (
	"csn/csn_tests_helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit(t *testing.T) {
	connectionString, pool, resource := csn_tests_helper.SetuMySql()
	defer csn_tests_helper.Purge(pool, resource)
	db, initErr := Init(*connectionString)
	require.NoError(t, initErr, "Could not init the database")
	dbErr := db.Ping()
	require.NoError(t, dbErr, "Failed to ping database")
}
