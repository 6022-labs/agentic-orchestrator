package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/settings"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const bridgeKindSlack = "slack"

// BridgeRegistrar calls the broker's upsert endpoint on startup with the
// configured bridge id. The broker returns 200 if the id was already known
// and 201 otherwise — both are success from our perspective. No local
// persistence needed: the id is owned by this slack-bridge's configuration.
type BridgeRegistrar struct {
	logger   *zap.Logger
	client   clients.ConversationBrokerClientInterface
	settings *settings.ConversationBrokerSettings
}

func NewBridgeRegistrar(
	logger *zap.Logger,
	client clients.ConversationBrokerClientInterface,
	settings *settings.ConversationBrokerSettings,
) *BridgeRegistrar {
	return &BridgeRegistrar{
		logger:   logger,
		client:   client,
		settings: settings,
	}
}

// Ensure registers (or confirms) this slack-bridge instance with the broker.
// Called at startup; broker is idempotent (200 if known, 201 otherwise).
func (r *BridgeRegistrar) Ensure() (uuid.UUID, error) {
	r.logger.Info("Upserting slack bridge with conversation broker",
		zap.String("bridgeId", r.settings.BridgeID.String()),
		zap.String("baseUrl", r.settings.CallbackUrl),
	)

	id := r.settings.BridgeID
	bridge, err := r.client.RegisterBridge(&models.ExternalRegisterBridgeRequest{
		ID:      &id,
		Kind:    bridgeKindSlack,
		BaseUrl: r.settings.CallbackUrl,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to upsert bridge with broker: %w", err)
	}
	if bridge.ID != r.settings.BridgeID {
		return uuid.Nil, fmt.Errorf("broker returned unexpected bridge id %s (expected %s)", bridge.ID, r.settings.BridgeID)
	}

	r.logger.Info("Slack bridge confirmed with conversation broker",
		zap.String("bridgeId", bridge.ID.String()),
	)
	return bridge.ID, nil
}
