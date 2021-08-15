package db

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

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
		db, err = gorm.Open("mysql", connStr)
		return err
	}

	_ = backoff.RetryNotify(connect, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
		fmt.Println(err.Error())
	})

	client.Gorm = db
	client.Gorm.DB().SetMaxOpenConns(1)
	client.Gorm.DB().SetMaxIdleConns(0)

	return &DBClient{
		client,
	}, nil
}

// CloseConnection closes a mysql connection
func (client *MySQLClient) CloseConnection() error {
	err := client.Gorm.Close()
	if err != nil {
		return err
	}

	return nil
}
