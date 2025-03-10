package main

import (
	"orders-service/helpers"
	"orders-service/httphandler"
	"orders-service/models"

	"github.com/IvanSkripnikov/go-gormdb"
	logger "github.com/IvanSkripnikov/go-logger"
)

func main() {
	logger.Debug("Service starting")

	// регистрация общих метрик
	helpers.RegisterCommonMetrics()

	// настройка всех конфигов
	config, err := models.LoadConfig()
	if err != nil {
		logger.Fatalf("Config error: %v", err)
	}

	// настройка коннекта к БД
	helpers.InitDatabase(config.Database)

	// выполнение миграций
	migrationModels := models.GetModels()
	gormdb.ApplyMigrationsForClient(models.ServiceDatabase, migrationModels...)

	// инициализация REST-api
	httphandler.InitHTTPServer()

	logger.Info("Service started")
}
