package db

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
		db, err = gorm.Open("postgres", connStr)
		return err
	}

	_ = backoff.RetryNotify(connect, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
		fmt.Println(err.Error())
	})

	client.Gorm = db
	client.Gorm.DB().SetMaxOpenConns(1)
	client.Gorm.DB().SetMaxIdleConns(0)
	client.Gorm.LogMode(false)

	return &DBClient{
		client,
	}, nil
}

// CloseConnection closes a postgres connection
func (client *PostgresClient) CloseConnection() error {
	err := client.Gorm.Close()
	if err != nil {
		return err
	}

	return nil
}
