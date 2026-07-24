package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"go.uber.org/zap"
)

const agentUrlCachePathTemplate = "agent-urls/%s"

// AgentUrlResolver caches the ENS-resolved node URL per domain (TTL enforced
// by the store). Failures are never cached.
type AgentUrlResolver struct {
	logger                        *zap.Logger
	cacheStorage                  interfaces.CacheStorageInterface
	ensUrlResolver                interfaces.EnsUrlResolverInterface
	remoteAgentResolutionSettings *settings.RemoteAgentResolutionSettings
	resolutionLocks               sync.Map
}

func NewAgentUrlResolver(
	logger *zap.Logger,
	cacheStorage interfaces.CacheStorageInterface,
	ensUrlResolver interfaces.EnsUrlResolverInterface,
	remoteAgentResolutionSettings *settings.RemoteAgentResolutionSettings,
) *AgentUrlResolver {
	return &AgentUrlResolver{
		logger:                        logger,
		cacheStorage:                  cacheStorage,
		ensUrlResolver:                ensUrlResolver,
		remoteAgentResolutionSettings: remoteAgentResolutionSettings,
	}
}

func (r *AgentUrlResolver) ResolveUrl(ensDomain string) (string, error) {
	ensDomain = strings.ToLower(strings.TrimSpace(ensDomain))
	cachePath := fmt.Sprintf(agentUrlCachePathTemplate, ensDomain)

	if url, ok := r.readCachedUrl(cachePath); ok {
		return url, nil
	}

	lock, _ := r.resolutionLocks.LoadOrStore(ensDomain, &sync.Mutex{})
	mutex := lock.(*sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	if url, ok := r.readCachedUrl(cachePath); ok {
		return url, nil
	}

	return r.resolveAndCacheUrl(ensDomain, cachePath)
}

func (r *AgentUrlResolver) readCachedUrl(cachePath string) (string, bool) {
	cached, err := r.cacheStorage.Read(cachePath)
	if err != nil || cached == "" {
		return "", false
	}
	return cached, true
}

func (r *AgentUrlResolver) resolveAndCacheUrl(ensDomain string, cachePath string) (string, error) {
	url, err := r.ensUrlResolver.ResolveUrl(ensDomain)
	if err != nil {
		return "", fmt.Errorf("failed to resolve url for %s: %w", ensDomain, err)
	}
	if url == "" {
		return "", fmt.Errorf("no url published for %s", ensDomain)
	}

	ttl := time.Duration(r.remoteAgentResolutionSettings.CacheTtlSeconds) * time.Second
	if err := r.cacheStorage.Store(cachePath, url, &ttl); err != nil {
		r.logger.Warn("Failed to cache agent url", zap.String("ensDomain", ensDomain), zap.Error(err))
	}

	return url, nil
}
