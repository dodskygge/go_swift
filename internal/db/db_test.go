package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "go_swift")

	database, err := ConnectDB()

	assert.NoError(t, err)
	assert.NotNil(t, database)

	if database != nil {
		database.Close()
	}
}
