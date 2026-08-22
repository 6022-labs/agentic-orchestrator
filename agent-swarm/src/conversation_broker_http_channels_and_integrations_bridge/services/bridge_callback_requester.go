package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/clients"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const jsonContentType = "application/json"

type BridgeCallbackRequester struct {
	logger        *zap.Logger
	clientFactory clients.ClientFactory
}

func NewBridgeCallbackRequester(logger *zap.Logger, clientFactory clients.ClientFactory) *BridgeCallbackRequester {
	return &BridgeCallbackRequester{
		logger:        logger,
		clientFactory: clientFactory,
	}
}

func (r *BridgeCallbackRequester) PostCompletion(
	ctx context.Context,
	bridge *models.Bridge,
	conversation *models.Conversation,
	agent *models.Participant,
	content string,
) error {
	if agent.CollectionAddress == nil || agent.CollectionTokenId == nil {
		return fmt.Errorf("participant is not an agent (missing on-chain identity); cannot send callback")
	}

	client, err := r.clientFactory(bridge.BaseUrl)
	if err != nil {
		return fmt.Errorf("failed to build bridge client: %w", err)
	}

	body, err := json.Marshal(clients.RequestsPostBrokerCompletionToSlackRequest{
		ConversationId:    lo.ToPtr(conversation.ID.String()),
		CollectionAddress: agent.CollectionAddress,
		CollectionTokenId: lo.ToPtr(int(*agent.CollectionTokenId)),
		Content:           &content,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal callback body: %w", err)
	}

	rsp, err := client.PostBrokerCallbackWithBodyWithResponse(
		ctx,
		jsonContentType,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	if rsp.StatusCode() < http.StatusOK || rsp.StatusCode() >= http.StatusMultipleChoices {
		r.logger.Warn(
			"Bridge callback returned non-2xx",
			zap.String("baseUrl", bridge.BaseUrl),
			zap.Int("status_code", rsp.StatusCode()),
			zap.String("body", string(rsp.Body)),
		)
		return fmt.Errorf("unexpected status code: %d", rsp.StatusCode())
	}

	return nil
}
