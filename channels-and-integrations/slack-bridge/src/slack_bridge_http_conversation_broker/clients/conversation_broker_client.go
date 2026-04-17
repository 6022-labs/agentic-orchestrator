package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/settings"
	"github.com/google/uuid"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	REGISTER_BRIDGE_ENDPOINT                    = "/bridges"
	REGISTER_AGENT_ENDPOINT                     = "/participants/agents"
	REGISTER_USER_ENDPOINT                      = "/participants/users"
	REGISTER_CONVERSATION_ENDPOINT              = "/conversations"
	GET_CONVERSATION_ENDPOINT                   = "/conversations/%s"
	UPDATE_CONVERSATION_DESCRIPTION_ENDPOINT    = "/conversations/%s/description"
	FORCE_CONCLUDE_ENDPOINT                     = "/conversations/%s/force-conclude"
	ATTACH_PARTICIPANT_ENDPOINT                 = "/conversations/%s/participants"
	DETACH_PARTICIPANT_ENDPOINT                 = "/conversations/%s/participants/%s"
	INGEST_MESSAGE_ENDPOINT                     = "/conversations/%s/messages"
)

type ConversationBrokerClientInterface interface {
	RegisterBridge(body *models.ExternalRegisterBridgeRequest) (*models.ExternalBridgeResponse, error)
	RegisterAgent(body *models.ExternalRegisterAgentRequest) (*models.ExternalParticipantResponse, error)
	RegisterUser(body *models.ExternalRegisterUserRequest) (*models.ExternalParticipantResponse, error)
	RegisterConversation(body *models.ExternalRegisterConversationRequest) (*models.ExternalConversationResponse, error)
	GetConversation(conversationId uuid.UUID) (*models.ExternalConversationResponse, error)
	UpdateConversationDescription(conversationId uuid.UUID, body *models.ExternalUpdateConversationDescriptionRequest) error
	ForceConcludeConversation(conversationId uuid.UUID) error
	AttachParticipant(conversationId uuid.UUID, body *models.ExternalConversationParticipantRequest) error
	DetachParticipant(conversationId uuid.UUID, participantId uuid.UUID) error
	IngestMessage(conversationId uuid.UUID, body *models.ExternalIngestMessageRequest) (*models.ExternalConversationMessageResponse, error)
}

type ConversationBrokerClient struct {
	logger                     *zap.Logger
	httpClient                 *http.Client
	conversationBrokerSettings *settings.ConversationBrokerSettings
}

type NewConversationBrokerClientParams struct {
	dig.In

	Logger                     *zap.Logger
	HttpClient                 *http.Client `name:"conversation-broker-http-client"`
	ConversationBrokerSettings *settings.ConversationBrokerSettings
}

func NewConversationBrokerClient(params NewConversationBrokerClientParams) *ConversationBrokerClient {
	return &ConversationBrokerClient{
		logger:                     params.Logger,
		httpClient:                 params.HttpClient,
		conversationBrokerSettings: params.ConversationBrokerSettings,
	}
}

func (c *ConversationBrokerClient) RegisterBridge(body *models.ExternalRegisterBridgeRequest) (*models.ExternalBridgeResponse, error) {
	var out models.ExternalBridgeResponse
	if err := c.doJSON(http.MethodPost, REGISTER_BRIDGE_ENDPOINT, body, &out, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}

// RegisterAgent is idempotent on the broker side: 200 if the agent was already
// known, 201 if newly created. Both are success.
func (c *ConversationBrokerClient) RegisterAgent(body *models.ExternalRegisterAgentRequest) (*models.ExternalParticipantResponse, error) {
	var out models.ExternalParticipantResponse
	if err := c.doJSON(http.MethodPost, REGISTER_AGENT_ENDPOINT, body, &out, http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ConversationBrokerClient) RegisterUser(body *models.ExternalRegisterUserRequest) (*models.ExternalParticipantResponse, error) {
	var out models.ExternalParticipantResponse
	if err := c.doJSON(http.MethodPost, REGISTER_USER_ENDPOINT, body, &out, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ConversationBrokerClient) RegisterConversation(body *models.ExternalRegisterConversationRequest) (*models.ExternalConversationResponse, error) {
	var out models.ExternalConversationResponse
	if err := c.doJSON(http.MethodPost, REGISTER_CONVERSATION_ENDPOINT, body, &out, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ConversationBrokerClient) GetConversation(conversationId uuid.UUID) (*models.ExternalConversationResponse, error) {
	var out models.ExternalConversationResponse
	endpoint := fmt.Sprintf(GET_CONVERSATION_ENDPOINT, conversationId)
	if err := c.doJSON(http.MethodGet, endpoint, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ConversationBrokerClient) UpdateConversationDescription(conversationId uuid.UUID, body *models.ExternalUpdateConversationDescriptionRequest) error {
	endpoint := fmt.Sprintf(UPDATE_CONVERSATION_DESCRIPTION_ENDPOINT, conversationId)
	return c.doJSON(http.MethodPatch, endpoint, body, nil, http.StatusNoContent)
}

func (c *ConversationBrokerClient) ForceConcludeConversation(conversationId uuid.UUID) error {
	endpoint := fmt.Sprintf(FORCE_CONCLUDE_ENDPOINT, conversationId)
	return c.doJSON(http.MethodPost, endpoint, nil, nil, http.StatusNoContent)
}

func (c *ConversationBrokerClient) AttachParticipant(conversationId uuid.UUID, body *models.ExternalConversationParticipantRequest) error {
	endpoint := fmt.Sprintf(ATTACH_PARTICIPANT_ENDPOINT, conversationId)
	return c.doJSON(http.MethodPost, endpoint, body, nil, http.StatusCreated)
}

func (c *ConversationBrokerClient) DetachParticipant(conversationId uuid.UUID, participantId uuid.UUID) error {
	endpoint := fmt.Sprintf(DETACH_PARTICIPANT_ENDPOINT, conversationId, participantId)
	return c.doJSON(http.MethodDelete, endpoint, nil, nil, http.StatusNoContent)
}

func (c *ConversationBrokerClient) IngestMessage(conversationId uuid.UUID, body *models.ExternalIngestMessageRequest) (*models.ExternalConversationMessageResponse, error) {
	var out models.ExternalConversationMessageResponse
	endpoint := fmt.Sprintf(INGEST_MESSAGE_ENDPOINT, conversationId)
	if err := c.doJSON(http.MethodPost, endpoint, body, &out, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ConversationBrokerClient) doJSON(method, endpoint string, body any, out any, expectedStatuses ...int) error {
	url := c.conversationBrokerSettings.BaseUrl + endpoint

	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if !slices.Contains(expectedStatuses, resp.StatusCode) {
		responseBody, _ := io.ReadAll(resp.Body)
		c.logger.Warn(
			"Unexpected status from conversation broker",
			zap.String("url", url),
			zap.String("method", method),
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(responseBody)),
		)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}
