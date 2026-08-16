package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/helpers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/settings"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// BridgeRegistrar calls the broker's upsert endpoint on startup with the
// configured bridge id. The broker returns 200 if the id was already known
// and 201 otherwise — both are success from our perspective. No local
// persistence needed: the id is owned by this slack-bridge's configuration.
type BridgeRegistrar struct {
	logger                     *zap.Logger
	client                     clients.ClientWithResponsesInterface
	conversationBrokerSettings *settings.ConversationBrokerSettings
}

func NewBridgeRegistrar(
	logger *zap.Logger,
	client clients.ClientWithResponsesInterface,
	conversationBrokerSettings *settings.ConversationBrokerSettings,
) *BridgeRegistrar {
	return &BridgeRegistrar{
		logger:                     logger,
		client:                     client,
		conversationBrokerSettings: conversationBrokerSettings,
	}
}

func (r *BridgeRegistrar) Ensure() (uuid.UUID, error) {
	r.logger.Info("Upserting slack bridge with conversation broker",
		zap.String("bridgeId", r.conversationBrokerSettings.BridgeID.String()),
		zap.String("baseUrl", r.conversationBrokerSettings.CallbackUrl),
	)

	body, err := json.Marshal(clients.RequestsRegisterBridgeRequest{
		Id:      lo.ToPtr(r.conversationBrokerSettings.BridgeID.String()),
		Kind:    lo.ToPtr(clients.BridgeKindSlack),
		BaseUrl: &r.conversationBrokerSettings.CallbackUrl,
	})
	if err != nil {
		return uuid.Nil, err
	}

	rsp, err := r.client.PostBridgesWithBodyWithResponse(context.Background(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to upsert bridge with broker: %w", err)
	}

	payload := rsp.JSON201
	if payload == nil {
		payload = rsp.JSON200
	}
	if payload == nil {
		return uuid.Nil, fmt.Errorf("failed to upsert bridge with broker: unexpected status code: %d, response: %s", rsp.StatusCode(), string(rsp.Body))
	}

	bridgeId := helpers.ParseUuid(payload.Id)
	if bridgeId != r.conversationBrokerSettings.BridgeID {
		return uuid.Nil, fmt.Errorf("broker returned unexpected bridge id %s (expected %s)", bridgeId, r.conversationBrokerSettings.BridgeID)
	}

	r.logger.Info("Slack bridge confirmed with conversation broker",
		zap.String("bridgeId", bridgeId.String()),
	)
	return bridgeId, nil
}
