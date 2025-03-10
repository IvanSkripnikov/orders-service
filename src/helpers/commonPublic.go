package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"orders-service/models"
	"time"

	logger "github.com/IvanSkripnikov/go-logger"
	"github.com/redis/go-redis/v9"
)

var Config *models.Config

func InitConfig(cfg *models.Config) {
	Config = cfg
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func FormatResponse(w http.ResponseWriter, httpStatus int, category string) {
	w.WriteHeader(httpStatus)

	data := ResponseData{
		"error": "Unsuccessfull request",
	}
	SendResponse(w, data, category, httpStatus)
}

func CreateQueryWithScalarResponse(method, url string, data any) (string, error) {
	var err error
	var response string

	jsonData, err := json.Marshal(data)
	if err != nil {
		return response, err
	}
	logger.Infof("json data: %v", string(jsonData))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return response, err
	}

	resp, err := client.Do(req)
	logger.Infof("response for make deposit: %v", resp.Body)
	if err != nil {
		return response, err
	}

	var result map[string]string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &result)

	logger.Infof("Data from create payment %v", result)

	// Преобразуем JSON-строку в map
	if err != nil {
		return response, err
	}

	response, ok := result["response"]
	if !ok {
		return "", errors.New("failed to create payment")
	}

	return response, nil
}

func SendNotification(message map[string]interface{}) {
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf("Error connection to Redis: %v", err)
	}

	_, err = redisClient.XAdd(context.Background(), &redis.XAddArgs{
		Stream: Config.Redis.Stream,
		Values: message,
	}).Result()
	if err != nil {
		logger.Fatalf("Error sending to redis stream: %v", err)
	} else {
		logger.Info("Succsessfuly send to stream")
	}
}
