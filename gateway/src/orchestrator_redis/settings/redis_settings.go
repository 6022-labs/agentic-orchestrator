package settings

import (
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type RedisSettings struct {
	DB       int64
	Host     string
	Password string
}

func NewRedisSettings(logger *zap.Logger) *RedisSettings {
	redisHost := os.Getenv("REDIS_HOST")
	if len(strings.TrimSpace(redisHost)) == 0 {
		logger.Fatal("please set your REDIS_HOST value in your environment")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if len(strings.TrimSpace(redisPassword)) == 0 {
		logger.Warn("redis password is not set")
	}

	strRedisDB := os.Getenv("REDIS_DB")
	if len(strings.TrimSpace(strRedisDB)) == 0 {
		logger.Warn("redis db is not set, using default value 0")
		strRedisDB = "0"
	}

	redisDB, err := strconv.ParseInt(strRedisDB, 10, 64)
	if err != nil {
		logger.Fatal("Error, your REDIS_DB value is not a valid number")
	}

	return &RedisSettings{
		DB:       redisDB,
		Host:     redisHost,
		Password: redisPassword,
	}
}
