package models

import (
	"os"

	"github.com/IvanSkripnikov/go-gormdb"
	"gorm.io/gorm/schema"
)

type Config struct {
	Database gormdb.Database
}

func LoadConfig() (*Config, error) {
	return &Config{
		Database: gormdb.Database{
			Address:  os.Getenv("DB_ADDRESS"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DB:       os.Getenv("DB_NAME"),
		},
	}, nil
}

func GetRequiredVariables() []string {
	return []string{
		// Обязательные переменные окружения для подключения к БД сервиса
		"DB_ADDRESS", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
	}
}

func GetModels() []schema.Tabler {
	return []schema.Tabler{
		Order{},
	}
}
