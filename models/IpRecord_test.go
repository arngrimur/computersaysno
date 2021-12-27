package models

import (
	"csn/db"
	"database/sql"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const ip = "1.1.1.1"

var connString *string
var pool *dockertest.Pool
var resource *dockertest.Resource
var sqlDb *sql.DB
var lookFor = IpRecord{
	ip: ip,
}

func TestMain(m *testing.M) {
	var testDbConfig = db.DbConfig{
		DbSecrets: db.DbSecrets{
			RootPassword: "secret",
			MysqlUser:    "testuser",
			MysqlPwd:     "testpassword",
		},
		HostConfig: db.HostConfig{
			AutoRemove:    true,
			RestartPolicy: "no",
		},
		ExpireTime:   uint(240),
		DatabaseName: "csn_db",
	}
	connString, pool, resource = db.SetupDatbase(testDbConfig)
	var err error
	sqlDb, err = db.InitDatabase(*connString)
	if err != nil {
		return
	}
	exitVal := m.Run()
	db.Purge(pool, resource)
	os.Exit(exitVal)
}

func TestIpRecord_New(t *testing.T) {
	record := NewIpRecord(ip)
	assert.Equal(t, ip, record.ip)
	assert.Equal(t, uint8(1), record.hitCount)
}

func TestIpRecord_CountIncrease(t *testing.T) {
	record := NewIpRecord(ip)
	record.IncreaseHitCount()
	assert.Equal(t, uint8(2), record.hitCount)
}

func TestIpRecord_CountIncreaseStops(t *testing.T) {
	record := NewIpRecord(ip)
	record.hitCount = 0xff
	record.IncreaseHitCount()
	assert.Equal(t, uint8(255), record.hitCount)
}

func TestIpRecord_CountIsNotBlocked(t *testing.T) {
	record := NewIpRecord(ip)
	record.hitCount = uint8(3)
	assert.Equal(t, false, record.IsBlocked(ip))
}

func TestIpRecord_CountIsBlocked(t *testing.T) {
	record := NewIpRecord(ip)
	record.hitCount = uint8(4)
	assert.Equal(t, true, record.IsBlocked(ip))
}

func TestIpRecord_GetIp(t *testing.T) {
	record := NewIpRecord(ip)
	assert.Equal(t, ip, record.GetIp())
}

func TestIpRecord_GetHitCount(t *testing.T) {
	record := NewIpRecord(ip)
	record.IncreaseHitCount()
	assert.Equal(t, uint8(2), record.GetHitCount())
}

func TestIpRecord_Create(t *testing.T) {
	record := NewIpRecord(ip)
	result, err := record.Create(sqlDb)
	require.NoError(t, err, "Failed to store IpRecord", err)
	affected, err := result.RowsAffected()
	require.NoError(t, err, "Rows affected error")
	assert.Equal(t, int64(1), affected)
}

func TestIpRecord_Read(t *testing.T) {

	recordRead, err := lookFor.Read(sqlDb)
	require.NoError(t, err, "Failed to read from database")
	assert.Equal(t, ip, recordRead.ip)
	assert.Equal(t, uint8(1), recordRead.hitCount)
}

func TestIpRecord_Update(t *testing.T) {
	recordRead, _ := lookFor.Read(sqlDb)
	recordRead.IncreaseHitCount()
	result, err := recordRead.Update(sqlDb)
	require.NoError(t, err, "Failed to update record in database")
	affected, err := result.RowsAffected()
	require.NoError(t, err, "No affected rows")
	assert.Equal(t, int64(1), affected)
	recordRead, _ = lookFor.Read(sqlDb)
	assert.Equal(t, ip, recordRead.ip)
	assert.Equal(t, uint8(2), recordRead.hitCount)
}

func TestIpRecord_Delete(t *testing.T) {
	record, err := lookFor.Delete(sqlDb)
	require.NoError(t, err, "Could no delete record")
	assert.Equal(t, ip, record.ip)
	assert.Equal(t, uint8(2), record.hitCount)
	_, err2 := lookFor.Read(sqlDb)
	require.Error(t, err2, "No IpRecord shall exist!")
}
