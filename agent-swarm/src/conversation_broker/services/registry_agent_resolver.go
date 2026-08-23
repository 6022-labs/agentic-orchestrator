package services

import (
	"context"

	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/traces/interfaces"
	"go.uber.org/zap"
)

const (
	registryAgentByCollectionCachePathTemplate = "registry-agents/collection/%s/%d"
	registryAgentByEnsDomainCachePathTemplate  = "registry-agents/ens/%s"
)

// RegistryAgentResolver caches the registry's agent projection under every lookup key, so one fetch warms both paths.
const registryAgentCacheName = "registry_agent"

type RegistryAgentResolver struct {
	logger                          *zap.Logger
	brokerMetrics                   metrics_interfaces.BrokerMetricsInterface
	brokerTracer                    traces_interfaces.BrokerTracerInterface
	cacheStorage                    interfaces.CacheStorageInterface
	agentRegistryRequester          interfaces.AgentRegistryRequesterInterface
	registryAgentResolutionSettings *settings.RegistryAgentResolutionSettings
	resolutionLocks                 sync.Map
}

func NewRegistryAgentResolver(
	logger *zap.Logger,
	brokerMetrics metrics_interfaces.BrokerMetricsInterface,
	brokerTracer traces_interfaces.BrokerTracerInterface,
	cacheStorage interfaces.CacheStorageInterface,
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	registryAgentResolutionSettings *settings.RegistryAgentResolutionSettings,
) *RegistryAgentResolver {
	return &RegistryAgentResolver{
		logger:                          logger,
		brokerMetrics:                   brokerMetrics,
		brokerTracer:                    brokerTracer,
		cacheStorage:                    cacheStorage,
		agentRegistryRequester:          agentRegistryRequester,
		registryAgentResolutionSettings: registryAgentResolutionSettings,
	}
}

func (r *RegistryAgentResolver) ResolveByCollection(
	ctx context.Context,
	collectionAddress string,
	collectionTokenId uint64,
) (*models.RegistryAgent, error) {
	return r.resolve(ctx,
		fmt.Sprintf("%s/%d", collectionAddress, collectionTokenId),
		r.collectionCachePath(collectionAddress, collectionTokenId),
		func() (*models.RegistryAgent, error) {
			return r.agentRegistryRequester.GetAgent(ctx, collectionAddress, collectionTokenId)
		},
	)
}

func (r *RegistryAgentResolver) ResolveByEnsDomain(
	ctx context.Context,
	ensDomain string,
) (*models.RegistryAgent, error) {
	ensDomain = strings.ToLower(strings.TrimSpace(ensDomain))

	return r.resolve(ctx, ensDomain, r.ensDomainCachePath(ensDomain), func() (*models.RegistryAgent, error) {
		return r.agentRegistryRequester.GetMintedAgentByEnsDomain(ctx, ensDomain)
	})
}

func (r *RegistryAgentResolver) collectionCachePath(collectionAddress string, collectionTokenId uint64) string {
	return fmt.Sprintf(
		registryAgentByCollectionCachePathTemplate,
		strings.ToLower(strings.TrimSpace(collectionAddress)),
		collectionTokenId,
	)
}

func (r *RegistryAgentResolver) ensDomainCachePath(ensDomain string) string {
	return fmt.Sprintf(registryAgentByEnsDomainCachePathTemplate, strings.ToLower(strings.TrimSpace(ensDomain)))
}

// resolve takes a per-key lock so concurrent misses collapse into one registry call.
func (r *RegistryAgentResolver) resolve(
	ctx context.Context,
	agentName string,
	cachePath string,
	fetch func() (*models.RegistryAgent, error),
) (resolved *models.RegistryAgent, err error) {
	ctx, span := r.brokerTracer.StartAgentResolution(ctx, metrics_interfaces.ResolverRegistry, agentName)
	defer func() {
		if err != nil {
			span.Fail(err)
		}
		r.brokerMetrics.RecordAgentResolution(ctx, metrics_interfaces.ResolverRegistry, err)
		span.End()
	}()

	if agent, ok := r.readCachedAgent(ctx, cachePath); ok {
		r.brokerMetrics.RecordCacheResult(ctx, registryAgentCacheName, true)
		return agent, nil
	}

	lock, _ := r.resolutionLocks.LoadOrStore(cachePath, &sync.Mutex{})
	mutex := lock.(*sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	if agent, ok := r.readCachedAgent(ctx, cachePath); ok {
		r.brokerMetrics.RecordCacheResult(ctx, registryAgentCacheName, true)
		return agent, nil
	}
	r.brokerMetrics.RecordCacheResult(ctx, registryAgentCacheName, false)

	agent, err := fetch()
	if err != nil {
		return nil, err
	}
	// A miss stays uncached so a newly minted agent is not hidden for a whole TTL.
	if agent == nil {
		return nil, nil
	}

	r.cacheAgent(ctx, cachePath, agent)
	return agent, nil
}

func (r *RegistryAgentResolver) readCachedAgent(ctx context.Context, cachePath string) (*models.RegistryAgent, bool) {
	cached, err := r.cacheStorage.Read(ctx, cachePath)
	if err != nil || cached == "" {
		return nil, false
	}

	var agent models.RegistryAgent
	if err := json.Unmarshal([]byte(cached), &agent); err != nil {
		return nil, false
	}
	return &agent, true
}

func (r *RegistryAgentResolver) cacheAgent(ctx context.Context, requestedPath string, agent *models.RegistryAgent) {
	serialized, err := json.Marshal(agent)
	if err != nil {
		return
	}

	ttl := time.Duration(r.registryAgentResolutionSettings.CacheTtlSeconds) * time.Second
	for _, cachePath := range r.cachePathsFor(requestedPath, agent) {
		if err := r.cacheStorage.Store(ctx, cachePath, string(serialized), &ttl); err != nil {
			r.logger.Warn("Failed to cache registry agent", zap.String("cachePath", cachePath), zap.Error(err))
		}
	}
}

// cachePathsFor skips a key the agent lacks the fields to build.
func (r *RegistryAgentResolver) cachePathsFor(requestedPath string, agent *models.RegistryAgent) []string {
	paths := []string{requestedPath}

	if agent.EnsDomain != "" {
		if path := r.ensDomainCachePath(agent.EnsDomain); path != requestedPath {
			paths = append(paths, path)
		}
	}
	if agent.CollectionAddress != "" && agent.CollectionTokenId != nil {
		if path := r.collectionCachePath(agent.CollectionAddress, *agent.CollectionTokenId); path != requestedPath {
			paths = append(paths, path)
		}
	}

	return paths
}
