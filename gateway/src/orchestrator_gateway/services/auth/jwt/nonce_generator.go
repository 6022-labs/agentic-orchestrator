package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/services/interfaces"
)

type NonceGenerator struct {
}

func NewNonceGenerator(redisStorage interfaces.CacheStorageInterface) *NonceGenerator {
	return &NonceGenerator{}
}

func (g *NonceGenerator) GenerateNonce() (*string, error) {
	nonce, err := generateRandomNonce(16)
	if err != nil {
		return nil, err
	}

	return &nonce, nil
}

func generateRandomNonce(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	timestamp := time.Now().Unix()
	return hex.EncodeToString(bytes) + "-" + strconv.FormatInt(timestamp, 10), nil
}
