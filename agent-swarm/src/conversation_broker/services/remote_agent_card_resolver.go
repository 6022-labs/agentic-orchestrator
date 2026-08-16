package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
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
	registryAgentResolver         interfaces.RegistryAgentResolverInterface
	remoteAgentCardRequester      interfaces.RemoteAgentCardRequesterInterface
	agentCardSignatureVerifier    interfaces.AgentCardSignatureVerifierInterface
	remoteAgentResolutionSettings *settings.RemoteAgentResolutionSettings
	resolutionLocks               sync.Map
}

func NewRemoteAgentCardResolver(
	logger *zap.Logger,
	cacheStorage interfaces.CacheStorageInterface,
	agentUrlResolver interfaces.AgentUrlResolverInterface,
	registryAgentResolver interfaces.RegistryAgentResolverInterface,
	remoteAgentCardRequester interfaces.RemoteAgentCardRequesterInterface,
	agentCardSignatureVerifier interfaces.AgentCardSignatureVerifierInterface,
	remoteAgentResolutionSettings *settings.RemoteAgentResolutionSettings,
) *RemoteAgentCardResolver {
	return &RemoteAgentCardResolver{
		logger:                        logger,
		cacheStorage:                  cacheStorage,
		agentUrlResolver:              agentUrlResolver,
		registryAgentResolver:         registryAgentResolver,
		remoteAgentCardRequester:      remoteAgentCardRequester,
		agentCardSignatureVerifier:    agentCardSignatureVerifier,
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

func (r *RemoteAgentCardResolver) resolveAndCacheCard(
	ensDomain string,
	cachePath string,
) (*models.RemoteAgentCard, error) {
	agentUrl, err := r.agentUrlResolver.ResolveUrl(ensDomain)
	if err != nil {
		return nil, err
	}

	card, raw, err := r.remoteAgentCardRequester.GetCard(agentUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch 6022 card for %s: %w", ensDomain, err)
	}
	if card == nil {
		return nil, fmt.Errorf("no 6022 card served for %s", ensDomain)
	}

	if err := r.verifySignatures(ensDomain, card, raw); err != nil {
		return nil, err
	}

	if serialized, err := json.Marshal(card); err == nil {
		ttl := time.Duration(r.remoteAgentResolutionSettings.CacheTtlSeconds) * time.Second
		if err := r.cacheStorage.Store(cachePath, string(serialized), &ttl); err != nil {
			r.logger.Warn("Failed to cache remote agent card", zap.String("ensDomain", ensDomain), zap.Error(err))
		}
	}

	return card, nil
}

// verifySignatures resolves the expected signer only when a claim was made, so unsigned cards cost nothing.
func (r *RemoteAgentCardResolver) verifySignatures(ensDomain string, card *models.RemoteAgentCard, raw []byte) error {
	if len(card.Signatures) == 0 {
		r.logger.Info("Remote agent card carries no signature",
			zap.String("ensDomain", ensDomain), zap.String("result", "unsigned"))
		return nil
	}

	expectedSigner, err := r.resolveExpectedSigner(ensDomain)
	if err != nil {
		r.logger.Warn("Failed to resolve the expected agent card signer",
			zap.String("ensDomain", ensDomain), zap.String("result", "anchor_unresolved"), zap.Error(err))
		return err
	}

	if err := r.agentCardSignatureVerifier.Verify(raw, card.Signatures, expectedSigner); err != nil {
		r.logger.Error("Remote agent card signature did not verify",
			zap.String("ensDomain", ensDomain), zap.String("result", "invalid_signature"), zap.Error(err))
		return fmt.Errorf("invalid 6022 card signature for %s: %w", ensDomain, err)
	}

	r.logger.Debug("Remote agent card signature verified",
		zap.String("ensDomain", ensDomain), zap.String("result", "ok"))

	return nil
}

// resolveExpectedSigner reads the wallet the registry holds for the agent; a card claiming a
// signature we cannot anchor is rejected rather than waved through.
func (r *RemoteAgentCardResolver) resolveExpectedSigner(ensDomain string) (string, error) {
	agent, err := r.registryAgentResolver.ResolveByEnsDomain(ensDomain)
	if err != nil {
		return "", fmt.Errorf("failed to resolve the expected signer for %s: %w", ensDomain, err)
	}
	if agent == nil {
		return "", fmt.Errorf("registry does not know %s", ensDomain)
	}

	signer := agent.EvmAddress()
	if signer == "" {
		return "", fmt.Errorf("registry holds no EVM address for %s", ensDomain)
	}
	return signer, nil
}
