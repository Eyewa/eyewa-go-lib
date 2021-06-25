package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// NewSQLiteClient creates a new sqlite client
func NewSQLiteClient() *SQLiteClient {
	return &SQLiteClient{
		Path: config.SQLite.Path,
	}
}

// NewSQLiteClientFromConfig creates a new sqlite client from manual configuration
func NewSQLiteClientFromConfig(config Config) *SQLiteClient {
	return &SQLiteClient{
		Path: config.SQLite.Path,
	}
}

// OpenConnection opens a sqlite connection
func (client *SQLiteClient) OpenConnection() (*DBClient, error) {
	db, err := gorm.Open("sqlite3", client.Path)
	if err != nil {
		return nil, err
	}

	client.gorm = db

	return &DBClient{
		client,
	}, nil

}

// CloseConnection closes a sqlite connection
func (client *SQLiteClient) CloseConnection() error {
	err := client.gorm.Close()
	if err != nil {
		return err
	}

	return nil
}
