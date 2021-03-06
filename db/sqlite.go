package db

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (*SQLiteClient) migrateDB() error {
	return nil
}

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
	var (
		db  *gorm.DB
		err error
	)

	connect := func() error {
		db, err = gorm.Open(sqlite.Open(client.Path), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent)})
		return err
	}

	_ = backoff.RetryNotify(connect, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	client.Gorm = db
	sql, err := client.Gorm.DB()
	if err != nil {
		return nil, err
	}

	sql.SetMaxOpenConns(1)
	sql.SetMaxIdleConns(0)

	return &DBClient{
		client,
	}, nil
}

// CloseConnection closes a sqlite connection
func (client *SQLiteClient) CloseConnection() error {
	if client.Gorm != nil {
		sql, err := client.Gorm.DB()
		if err != nil {
			return err
		}

		err = sql.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
