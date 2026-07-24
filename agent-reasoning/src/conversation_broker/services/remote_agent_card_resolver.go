package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"go.uber.org/zap"
)

const remoteAgentCardCachePathTemplate = "remote-agent-cards/%s"

// RemoteAgentCardResolver composes the cached url resolution + card fetch,
// caching the card per ENS domain (TTL enforced by the store). Failures are
// never cached.
type RemoteAgentCardResolver struct {
	logger                        *zap.Logger
	cacheStorage                  interfaces.CacheStorageInterface
	agentUrlResolver              interfaces.AgentUrlResolverInterface
	remoteAgentCardRequester      interfaces.RemoteAgentCardRequesterInterface
	remoteAgentResolutionSettings *settings.RemoteAgentResolutionSettings
	resolutionLocks               sync.Map
}

func NewRemoteAgentCardResolver(
	logger *zap.Logger,
	cacheStorage interfaces.CacheStorageInterface,
	agentUrlResolver interfaces.AgentUrlResolverInterface,
	remoteAgentCardRequester interfaces.RemoteAgentCardRequesterInterface,
	remoteAgentResolutionSettings *settings.RemoteAgentResolutionSettings,
) *RemoteAgentCardResolver {
	return &RemoteAgentCardResolver{
		logger:                        logger,
		cacheStorage:                  cacheStorage,
		agentUrlResolver:              agentUrlResolver,
		remoteAgentCardRequester:      remoteAgentCardRequester,
		remoteAgentResolutionSettings: remoteAgentResolutionSettings,
	}
}

func (r *RemoteAgentCardResolver) ResolveCard(ensDomain string) (*models.RemoteAgentCard, error) {
	ensDomain = strings.ToLower(strings.TrimSpace(ensDomain))
	cachePath := fmt.Sprintf(remoteAgentCardCachePathTemplate, ensDomain)

	if card, ok := r.readCachedCard(cachePath); ok {
		return card, nil
	}

	lock, _ := r.resolutionLocks.LoadOrStore(ensDomain, &sync.Mutex{})
	mutex := lock.(*sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	if card, ok := r.readCachedCard(cachePath); ok {
		return card, nil
	}

	return r.resolveAndCacheCard(ensDomain, cachePath)
}

func (r *RemoteAgentCardResolver) readCachedCard(cachePath string) (*models.RemoteAgentCard, bool) {
	cached, err := r.cacheStorage.Read(cachePath)
	if err != nil {
		return nil, false
	}

	var card models.RemoteAgentCard
	if err := json.Unmarshal([]byte(cached), &card); err != nil || card.Protocol == "" {
		return nil, false
	}
	return &card, true
}

func (r *RemoteAgentCardResolver) resolveAndCacheCard(ensDomain string, cachePath string) (*models.RemoteAgentCard, error) {
	agentUrl, err := r.agentUrlResolver.ResolveUrl(ensDomain)
	if err != nil {
		return nil, err
	}

	card, err := r.remoteAgentCardRequester.GetCard(agentUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch 6022 card for %s: %w", ensDomain, err)
	}

	if serialized, err := json.Marshal(card); err == nil {
		ttl := time.Duration(r.remoteAgentResolutionSettings.CacheTtlSeconds) * time.Second
		if err := r.cacheStorage.Store(cachePath, string(serialized), &ttl); err != nil {
			r.logger.Warn("Failed to cache remote agent card", zap.String("ensDomain", ensDomain), zap.Error(err))
		}
	}

	return card, nil
}
