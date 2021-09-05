package db

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresClient creates a new postgres client
func NewPostgresClient() *PostgresClient {
	return &PostgresClient{
		nil,
		RDMS{
			Name:     config.Database.Name,
			Host:     config.Database.Host,
			Port:     config.Database.Port,
			User:     config.Database.User,
			Password: config.Database.Password,
			SSLMode:  config.Database.SSLMode,
		},
	}
}

// NewPostgresClientFromConfig creates a new postgres client from a manual configuration
func NewPostgresClientFromConfig(config Config) *PostgresClient {
	return &PostgresClient{
		nil,
		RDMS{
			Name:     config.Database.Name,
			Host:     config.Database.Host,
			Port:     config.Database.Port,
			User:     config.Database.User,
			Password: config.Database.Password,
			SSLMode:  config.Database.SSLMode,
		},
	}
}

// OpenConnection opens connection to postgres
func (client *PostgresClient) OpenConnection() (*DBClient, error) {
	var (
		db  *gorm.DB
		err error
	)

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s  dbname=%s",
		client.User,
		client.Password,
		client.Host,
		client.Port,
		client.Name,
	)

	if client.SSLMode == "disable" {
		connStr = fmt.Sprintf("%s sslmode=%s", connStr, client.SSLMode)
	}

	connect := func() error {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
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

// CloseConnection closes a postgres connection
func (client *PostgresClient) CloseConnection() error {
	sql, err := client.Gorm.DB()
	if err != nil {
		return err
	}

	err = sql.Close()
	if err != nil {
		return err
	}

	return nil
}
