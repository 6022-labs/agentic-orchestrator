package services

import (
	"context"

	"fmt"
	"strings"
	"sync"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/traces/interfaces"
	"go.uber.org/zap"
)

const agentUrlCacheName = "agent_url"

const agentUrlCachePathTemplate = "agent-urls/%s"

// AgentUrlResolver caches the ENS-resolved node URL per domain (TTL enforced
// by the store). Failures are never cached.
type AgentUrlResolver struct {
	logger                        *zap.Logger
	brokerMetrics                 metrics_interfaces.BrokerMetricsInterface
	brokerTracer                  traces_interfaces.BrokerTracerInterface
	cacheStorage                  interfaces.CacheStorageInterface
	ensUrlResolver                interfaces.EnsUrlResolverInterface
	remoteAgentResolutionSettings *settings.RemoteAgentResolutionSettings
	resolutionLocks               sync.Map
}

func NewAgentUrlResolver(
	logger *zap.Logger,
	brokerMetrics metrics_interfaces.BrokerMetricsInterface,
	brokerTracer traces_interfaces.BrokerTracerInterface,
	cacheStorage interfaces.CacheStorageInterface,
	ensUrlResolver interfaces.EnsUrlResolverInterface,
	remoteAgentResolutionSettings *settings.RemoteAgentResolutionSettings,
) *AgentUrlResolver {
	return &AgentUrlResolver{
		logger:                        logger,
		brokerMetrics:                 brokerMetrics,
		brokerTracer:                  brokerTracer,
		cacheStorage:                  cacheStorage,
		ensUrlResolver:                ensUrlResolver,
		remoteAgentResolutionSettings: remoteAgentResolutionSettings,
	}
}

func (r *AgentUrlResolver) ResolveUrl(ctx context.Context, ensDomain string) (resolved string, err error) {
	ensDomain = strings.ToLower(strings.TrimSpace(ensDomain))
	cachePath := fmt.Sprintf(agentUrlCachePathTemplate, ensDomain)

	ctx, span := r.brokerTracer.StartAgentResolution(ctx, metrics_interfaces.ResolverEns, ensDomain)
	defer func() {
		if err != nil {
			span.Fail(err)
		}
		r.brokerMetrics.RecordAgentResolution(ctx, metrics_interfaces.ResolverEns, err)
		span.End()
	}()

	if url, ok := r.readCachedUrl(ctx, cachePath); ok {
		r.brokerMetrics.RecordCacheResult(ctx, agentUrlCacheName, true)
		return url, nil
	}

	lock, _ := r.resolutionLocks.LoadOrStore(ensDomain, &sync.Mutex{})
	mutex := lock.(*sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	if url, ok := r.readCachedUrl(ctx, cachePath); ok {
		r.brokerMetrics.RecordCacheResult(ctx, agentUrlCacheName, true)
		return url, nil
	}
	r.brokerMetrics.RecordCacheResult(ctx, agentUrlCacheName, false)

	return r.resolveAndCacheUrl(ctx, ensDomain, cachePath)
}

func (r *AgentUrlResolver) readCachedUrl(ctx context.Context, cachePath string) (string, bool) {
	cached, err := r.cacheStorage.Read(ctx, cachePath)
	if err != nil || cached == "" {
		return "", false
	}
	return cached, true
}

func (r *AgentUrlResolver) resolveAndCacheUrl(ctx context.Context, ensDomain string, cachePath string) (string, error) {
	url, err := r.ensUrlResolver.ResolveUrl(ctx, ensDomain)
	if err != nil {
		return "", fmt.Errorf("failed to resolve url for %s: %w", ensDomain, err)
	}
	if url == "" {
		return "", fmt.Errorf("no url published for %s", ensDomain)
	}

	ttl := time.Duration(r.remoteAgentResolutionSettings.CacheTtlSeconds) * time.Second
	if err := r.cacheStorage.Store(ctx, cachePath, url, &ttl); err != nil {
		r.logger.Warn("Failed to cache agent url", zap.String("ensDomain", ensDomain), zap.Error(err))
	}

	return url, nil
}
