package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// NewSQLiteClient create a new sqlite client
func NewSQLiteClient() *SQLiteClient {
	return &SQLiteClient{
		Path: config.SQLite.Path,
	}
}

func (client *SQLiteClient) openConnection() (*DBClient, error) {
	db, err := gorm.Open("sqlite3", client.Path)
	if err != nil {
		return nil, err
	}

	client.gorm = db

	return &DBClient{
		client,
	}, nil

}

func (client *SQLiteClient) closeConnection() error {
	err := client.gorm.Close()
	if err != nil {
		return err
	}

	return nil
}
