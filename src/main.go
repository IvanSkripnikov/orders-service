package main

import (
	"authenticator/helpers"
	"authenticator/httphandler"
	"authenticator/logger"
	"authenticator/models"
	"fmt"
)

func main() {
	logger.Debug("Service starting")

	// регистрация общих метрик
	helpers.RegisterCommonMetrics()

	// настройка всех конфигов
	config, err := models.LoadConfig()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Config error: %v", err))
	}

	// настройка коннекта к БД
	_, err = helpers.InitDataBase(config.Database)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Cant initialize DB: %v", err))
	}

	// инициализация сессий
	helpers.SessionsMap = map[string]models.User{}

	// инициализация REST-api
	httphandler.InitHTTPServer()

	logger.Info("Service started")
}
