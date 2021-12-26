package csn_tests_helper

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestStart(t *testing.T) {
	connString, pool, resource := SetuMySql()
	resource.GetPort("3306")
	defer Purge(pool, resource)
	assert.Equal(t, "testuser:testpassword@(localhost:55124)/csn_db?parseTime=true", connString)
}