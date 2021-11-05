package db

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (client *PostgresClient) migrateDB() error {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		client.User,
		client.Password,
		client.Host,
		client.Port,
		"postgres",
	)

	// connect to the postgres db just to be able to run the create db statement
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return err
	}

	// check if db exists
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", client.Name)
	rs := db.Raw(stmt)
	if rs.Error != nil {
		return rs.Error
	}

	// if not create it
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		stmt := fmt.Sprintf("CREATE DATABASE %s;", client.Name)
		if rs := db.Exec(stmt); rs.Error != nil {
			return rs.Error
		}

		// close db connection
		sql, err := db.DB()
		defer func() {
			_ = sql.Close()
		}()
		if err != nil {
			return err
		}
	}

	return nil
}

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
	// migrate db if not exists
	if err := client.migrateDB(); err != nil {
		return nil, err
	}

	var (
		db  *gorm.DB
		err error
	)

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
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

	sql.SetMaxIdleConns(5)
	sql.SetConnMaxIdleTime(30 * time.Minute)
	sql.SetMaxOpenConns(20)

	return &DBClient{
		client,
	}, nil
}

// CloseConnection closes a postgres connection
func (client *PostgresClient) CloseConnection() error {
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
