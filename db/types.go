package db

import (
	"gorm.io/gorm"
)

// DBDriver supported db driver
type DBDriver string

// Config for all db env vars
type Config struct {
	Driver   string       `mapstructure:"db_driver"`
	Database RDMS         `mapstructure:",squash"`
	SQLite   SQLiteClient `mapstructure:",squash,omitempty"`
}

// DatabaseDriver interface for all db clients
type DatabaseDriver interface {
	OpenConnection() (*DBClient, error)
	CloseConnection() error

	migrateDB() error
}

// RDMS definition for general RDMS
type RDMS struct {
	Host     string `mapstructure:"db_host"`
	Port     string `mapstructure:"db_port"`
	Name     string `mapstructure:"db_database"`
	User     string `mapstructure:"db_user"`
	Password string `mapstructure:"db_password"`
	SSLMode  string `mapstructure:"db_ssl_mode"`
}

// SQLite sqlite client definition
type SQLiteClient struct {
	Gorm *gorm.DB
	Path string `mapstructure:"db_path"`
}

// MySQL mysql client definition
type MySQLClient struct {
	Gorm *gorm.DB
	RDMS
}

// PostgresClient postgres client definition
type PostgresClient struct {
	Gorm *gorm.DB
	RDMS
}

type DBClient struct {
	DatabaseDriver
}
