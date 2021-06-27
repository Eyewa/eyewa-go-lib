package db

import (
	"os"
	"testing"

	libErrs "github.com/eyewa/eyewa-go-lib/errors"
	"github.com/stretchr/testify/assert"
)

var vars map[string]string

func TestInitConfigWithPostgres(t *testing.T) {
	os.Clearenv()
	config = *new(Config)

	vars = map[string]string{
		"DB_DRIVER":   "postgres",
		"DB_HOST":     "localhost",
		"DB_USER":     "admin",
		"DB_PORT":     "5432",
		"DB_PASSWORD": "secret",
		"DB_DATABASE": "catalogconsumer",
		"DB_SSL_MODE": "disabled",
	}
	for e, v := range vars {
		os.Setenv(e, v)
	}

	_, _ = initConfig()
	assert.NotZero(t, config)

	assert.Equal(t, string(Postgres), config.Driver)
	assert.Equal(t, "catalogconsumer", config.Database.Name)
	assert.Equal(t, "admin", config.Database.User)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "secret", config.Database.Password)
	assert.Equal(t, *new(SQLiteClient), config.SQLite)
}

func TestMySQLClient(t *testing.T) {
	os.Clearenv()
	config = *new(Config)

	vars = map[string]string{
		"DB_DRIVER":   "mysql",
		"DB_HOST":     "localhost",
		"DB_USER":     "admin",
		"DB_PORT":     "3306",
		"DB_PASSWORD": "secret",
		"DB_DATABASE": "catalogconsumer",
	}
	for e, v := range vars {
		os.Setenv(e, v)
	}

	_, _ = initConfig()

	assert.Equal(t, string(MySQL), config.Driver)
	assert.Equal(t, "catalogconsumer", config.Database.Name)
	assert.Equal(t, "admin", config.Database.User)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "secret", config.Database.Password)
	assert.Equal(t, *new(SQLiteClient), config.SQLite)
}

func TestInitConfigWithSQLite(t *testing.T) {
	os.Clearenv()
	config = *new(Config)

	vars = map[string]string{
		"DB_DRIVER": "sqlite",
		"DB_PATH":   ":memory:",
	}
	for e, v := range vars {
		os.Setenv(e, v)
	}

	_, _ = initConfig()
	assert.NotZero(t, config)

	assert.Equal(t, string(SQLite), config.Driver)
	assert.Empty(t, config.Database.Name)
	assert.Empty(t, config.Database.User)
	assert.Empty(t, config.Database.Host)
	assert.Empty(t, config.Database.Password)
	assert.NotEqual(t, *new(SQLiteClient), config.SQLite)
	assert.Equal(t, ":memory:", config.SQLite.Path)

	client := NewSQLiteClient()
	_, err := client.openConnection()
	assert.Nil(t, err)

	err = client.closeConnection()
	assert.Nil(t, err)
}

func TestConnection(t *testing.T) {
	os.Clearenv()
	client = nil

	vars = map[string]string{
		"DB_DRIVER": "sqlite",
		"DB_PATH":   ":memory:",
	}
	for e, v := range vars {
		os.Setenv(e, v)
	}

	dbclient, err := OpenConnection()
	assert.Nil(t, err)
	assert.NotNil(t, dbclient)

	err = CloseConnection()
	assert.Nil(t, err)
}

func TestConnectionFail(t *testing.T) {
	os.Clearenv()
	client = nil

	vars = map[string]string{
		"DB_DRIVER":   "mssql",
		"DB_HOST":     "localhost",
		"DB_USER":     "admin",
		"DB_PORT":     "3306",
		"DB_PASSWORD": "secret",
		"DB_DATABASE": "catalogconsumer",
	}
	for e, v := range vars {
		os.Setenv(e, v)
	}

	dbclient, err := OpenConnection()
	assert.Nil(t, dbclient)
	assert.EqualError(t, libErrs.ErrorUnsupportedDBDriverSpecified, err.Error())

	err = CloseConnection()
	assert.EqualError(t, libErrs.ErrorNoDBClientFound, err.Error())
}
