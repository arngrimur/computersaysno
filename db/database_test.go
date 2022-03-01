package db

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	dir          = "/tmp/db"
	usernameFile = dir + "/username"
	passwordFile = dir + "/password"
)

func TestMain(m *testing.M) {
	defer tearDown()
	setup()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setup() {
	os.Mkdir(dir, 0777)
	_ = os.WriteFile(usernameFile, []byte("testuser"), 0400)
	_ = os.WriteFile(passwordFile, []byte("testpassword"), 0400)
}

func tearDown() {
	_ = os.RemoveAll(dir)
}

func TestGetConnectionsString(t *testing.T) {
	res := createConnectionString()
	assert.Equal(t, "jdbc:postgresql://testuser:testpassword@localhost:55144/csn_db?sslmode=disable", res)
}
