package services

import (
	"context"

	"encoding/json"
	"fmt"
	"strings"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"go.uber.org/zap"
)

const x402PaymentChallengeCacheName = "x402_payment_challenge"

const x402PaymentChallengeCachePathTemplate = "x402-payment-challenges/%s"

// X402PaymentChallengeCache caches the last 402 challenge per ENS domain; stale entries self-heal via a fresh 402.
type X402PaymentChallengeCache struct {
	logger              *zap.Logger
	brokerMetrics       metrics_interfaces.BrokerMetricsInterface
	cacheStorage        interfaces.CacheStorageInterface
	x402PaymentSettings *settings.X402PaymentSettings
}

func NewX402PaymentChallengeCache(
	logger *zap.Logger,
	brokerMetrics metrics_interfaces.BrokerMetricsInterface,
	cacheStorage interfaces.CacheStorageInterface,
	x402PaymentSettings *settings.X402PaymentSettings,
) *X402PaymentChallengeCache {
	return &X402PaymentChallengeCache{
		logger:              logger,
		brokerMetrics:       brokerMetrics,
		cacheStorage:        cacheStorage,
		x402PaymentSettings: x402PaymentSettings,
	}
}

func (c *X402PaymentChallengeCache) Read(ctx context.Context, ensDomain string) (*models.X402PaymentChallenge, bool) {
	cached, err := c.cacheStorage.Read(ctx, c.cachePath(ensDomain))
	if err != nil {
		// A cold cache surfaces here as an error, so this is the ordinary miss.
		c.brokerMetrics.RecordCacheResult(ctx, x402PaymentChallengeCacheName, false)
		return nil, false
	}

	var paymentChallenge models.X402PaymentChallenge
	if err := json.Unmarshal([]byte(cached), &paymentChallenge); err != nil || paymentChallenge.X402Version == 0 {
		c.brokerMetrics.RecordCacheResult(ctx, x402PaymentChallengeCacheName, false)
		return nil, false
	}

	c.brokerMetrics.RecordCacheResult(ctx, x402PaymentChallengeCacheName, true)
	return &paymentChallenge, true
}

func (c *X402PaymentChallengeCache) Store(
	ctx context.Context,
	ensDomain string,
	paymentChallenge *models.X402PaymentChallenge,
) {
	serialized, err := json.Marshal(paymentChallenge)
	if err != nil {
		c.logger.Warn(
			"Failed to serialize x402 payment requirements",
			zap.String("ensDomain", ensDomain),
			zap.Error(err),
		)
		return
	}

	ttl := time.Duration(c.x402PaymentSettings.ChallengeCacheTtlSeconds) * time.Second
	if err := c.cacheStorage.Store(ctx, c.cachePath(ensDomain), string(serialized), &ttl); err != nil {
		c.logger.Warn("Failed to cache x402 payment requirements", zap.String("ensDomain", ensDomain), zap.Error(err))
	}
}

func (c *X402PaymentChallengeCache) cachePath(ensDomain string) string {
	return fmt.Sprintf(x402PaymentChallengeCachePathTemplate, strings.ToLower(strings.TrimSpace(ensDomain)))
}
