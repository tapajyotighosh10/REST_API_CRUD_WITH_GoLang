package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName  string
	SSLMode string
}

func Conn(config *Config) (*gorm.DB, error) {
	dsn :=
		fmt.Sprintf(
			"host=%s  port=%s user=%s  password=%s  dbname=%s  sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
		)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, nil
	}
	return db, nil
}
