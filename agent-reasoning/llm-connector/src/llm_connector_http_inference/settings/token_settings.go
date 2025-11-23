package settings

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

type TokenSettings struct {
	MaxTokensPerRequest uint
}

func NewTokenSettings(
	logger *zap.Logger,
) *TokenSettings {
	maxTokensPerRequestStr := os.Getenv("MAX_TOKENS_PER_REQUEST")
	if len(strings.TrimSpace(maxTokensPerRequestStr)) == 0 {
		defaultMaxTokens := "100000"
		logger.Info("MAX_TOKENS_PER_REQUEST is not set, using default settings", zap.String("default", defaultMaxTokens))
		maxTokensPerRequestStr = defaultMaxTokens
	}

	var maxTokensPerRequest uint
	_, err := fmt.Sscanf(maxTokensPerRequestStr, "%d", &maxTokensPerRequest)
	if err != nil {
		logger.Warn("Invalid MAX_TOKENS_PER_REQUEST value, using default of 100000", zap.String("value", maxTokensPerRequestStr))
		maxTokensPerRequest = 100000
	}

	return &TokenSettings{
		MaxTokensPerRequest: maxTokensPerRequest,
	}
}
