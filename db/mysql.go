package db

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (client *MySQLClient) migrateDB() error {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		client.User,
		client.Password,
		client.Host,
		client.Port,
		"information_schema")

	// connect to the information_schema db - just to be able to run the create db statement
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return err
	}

	// check if db exists (if not create it)
	rs := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", client.Name))
	if rs.Error != nil {
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

	return nil
}

// NewMySQLClient create a new mysql client
func NewMySQLClient() *MySQLClient {
	return &MySQLClient{
		nil,
		RDMS{
			Name:     config.Database.Name,
			Host:     config.Database.Host,
			Port:     config.Database.Port,
			User:     config.Database.User,
			Password: config.Database.Password,
		},
	}
}

// NewMySQLClientFromConfig create a new mysql client from manual configuration
func NewMySQLClientFromConfig(config Config) *MySQLClient {
	return &MySQLClient{
		nil,
		RDMS{
			Name:     config.Database.Name,
			Host:     config.Database.Host,
			Port:     config.Database.Port,
			User:     config.Database.User,
			Password: config.Database.Password,
		},
	}
}

// OpenConnection opens connection to mysql
func (client *MySQLClient) OpenConnection() (*DBClient, error) {
	// migrate db if not exists
	if err := client.migrateDB(); err != nil {
		return nil, err
	}

	var (
		db  *gorm.DB
		err error
	)

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		client.User,
		client.Password,
		client.Host,
		client.Port,
		client.Name)

	connect := func() error {
		db, err = gorm.Open(mysql.Open(connStr), &gorm.Config{
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

	sql.SetMaxIdleConns(50)
	sql.SetConnMaxIdleTime(1 * time.Hour)
	sql.SetConnMaxLifetime(1 * time.Hour)
	sql.SetMaxOpenConns(50)

	return &DBClient{
		client,
	}, nil
}

// CloseConnection closes a mysql connection
func (client *MySQLClient) CloseConnection() error {
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
