package db

import (
	"fmt"

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

func (client *MySQLClient) openConnection() (*DBClient, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		client.User,
		client.Password,
		client.Host,
		client.Port,
		client.Name)

	db, err := gorm.Open("mysql", connStr)
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

func (client *MySQLClient) closeConnection() error {
	err := client.gorm.Close()
	if err != nil {
		return err
	}

	return nil
}
