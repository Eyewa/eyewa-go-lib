package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewPostgresClient create a new postgres client
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

func (client *PostgresClient) openConnection() (*DBClient, error) {
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

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	client.gorm = db
	client.gorm.DB().SetMaxOpenConns(1)
	client.gorm.DB().SetMaxIdleConns(0)

	return &DBClient{
		client,
	}, nil

}

func (client *PostgresClient) closeConnection() error {
	err := client.gorm.Close()
	if err != nil {
		return err
	}

	return nil
}
