package helpers

import (
	"context"
	"net"
	"orders-service/models"
	"strconv"

	logger "github.com/IvanSkripnikov/go-logger"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// Init Инициализация подключения к Redis.
func InitRedis(ctx context.Context, config models.Redis) {
	if _, err := strconv.Atoi(config.Port); err != nil {
		logger.Fatalf("Failed to parse on Redis port. Error: %v", err)
	}

	address := net.JoinHostPort(config.Address, config.Port)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: config.Password,
		DB:       config.DB,
	})
}
