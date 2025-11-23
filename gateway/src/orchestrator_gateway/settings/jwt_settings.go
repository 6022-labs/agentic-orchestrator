package settings

import (
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type JwtSettings struct {
	Hash       string
	LivingTime int
}

func NewJwtSettings(
	logger *zap.Logger,
) *JwtSettings {
	jwtHash := os.Getenv("JWT_HASH")
	jwtLivingTimeStr := os.Getenv("JWT_LIVING_TIME")

	if len(strings.TrimSpace(jwtHash)) == 0 {
		logger.Fatal("Error, please set your JWT_HASH value in your environment")
	}

	if len(strings.TrimSpace(jwtLivingTimeStr)) == 0 {
		logger.Fatal("Error, please set your JWT_LIVING_TIME value in your environment")
	}

	jwtLivingTime, err := strconv.Atoi(jwtLivingTimeStr)

	if err != nil {
		logger.Fatal("Error, JWT_LIVING_TIME is not a number.")
	}

	return &JwtSettings{
		Hash:       jwtHash,
		LivingTime: jwtLivingTime,
	}
}
