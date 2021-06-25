package db

import (
	"strings"

	libErrs "github.com/eyewa/eyewa-go-lib/errors"
	"github.com/ory/viper"
)

const (
	MySQL    DBDriver = "mysql"
	Postgres DBDriver = "postgres"
	SQLite   DBDriver = "sqlite"
)

var (
	client  *DBClient
	config  Config
	envVars = []string{
		"DB_DRIVER",
		"DB_HOST",
		"DB_USER",
		"DB_PORT",
		"DB_DATABASE",
		"DB_PASSWORD",
		"DB_PATH",
		"DB_SSL_MODE",
	}
)

func initConfig() (Config, error) {
	config = *new(Config)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for _, v := range envVars {
		if err := viper.BindEnv(v); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}

// OpenConnection opens a new db connection for chosen client
func OpenConnection() (*DBClient, error) {
	var err error

	config, err = initConfig()
	if err != nil {
		return nil, err
	}

	if config.Driver == "" {
		return nil, libErrs.ErrorNoDBDriverSpecified
	}

	switch strings.ToLower(config.Driver) {
	case string(MySQL):
		client = &DBClient{NewMySQLClient()}
	case string(Postgres):
		client = &DBClient{NewPostgresClient()}
	case string(SQLite):
		client = &DBClient{NewSQLiteClient()}
	}

	if client != nil {
		return client.OpenConnection()
	}

	return nil, libErrs.ErrorUnsupportedDBDriverSpecified
}

// CloseConnection close connection on db client
func CloseConnection() error {
	if client != nil {
		return client.CloseConnection()
	}

	return libErrs.ErrorNoDBClientFound
}
