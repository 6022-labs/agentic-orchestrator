package clients

import (
	"context"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/models"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type SlackClientInterface interface {
	PostMessage(
		ctx context.Context,
		request *models.ExternalPostMessageRequest,
	) (*models.ExternalPostMessageResponse, error)
	JoinConversation(
		ctx context.Context,
		request *models.ExternalJoinConversationRequest,
	) (*models.ExternalJoinConversationResponse, error)
	GetConversationInfo(
		ctx context.Context,
		request *models.ExternalGetConversationInfoRequest,
	) (*models.ExternalGetConversationInfoResponse, error)
	GetConversationMembers(
		ctx context.Context,
		request *models.ExternalGetConversationMembersRequest,
	) (*models.ExternalGetConversationMembersResponse, error)
	GetUsersInfo(
		ctx context.Context,
		request *models.ExternalGetUsersInfoRequest,
	) (*models.ExternalGetUsersInfoResponse, error)
}

const (
	POST_MESSAGE_ENDPOINT             = "/api/chat.postMessage"
	JOIN_CONVERSATION_ENDPOINT        = "/api/conversations.join"
	GET_CONVERSATION_INFO_ENDPOINT    = "/api/conversations.info"
	GET_CONVERSATION_MEMBERS_ENDPOINT = "/api/conversations.members"
	GET_USERS_INFO_ENDPOINT           = "/api/users.info"
)

type SlackClient struct {
	logger        *zap.Logger
	client        *http.Client
	slackSettings *settings.SlackSettings
}

type newSlackClientParams struct {
	dig.In

	Logger        *zap.Logger
	Client        *http.Client `name:"slack-http-client"`
	SlackSettings *settings.SlackSettings
}

func NewSlackClient(
	params newSlackClientParams,
) *SlackClient {
	return &SlackClient{
		logger:        params.Logger,
		client:        params.Client,
		slackSettings: params.SlackSettings,
	}
}

func (s *SlackClient) PostMessage(
	ctx context.Context,
	request *models.ExternalPostMessageRequest,
) (*models.ExternalPostMessageResponse, error) {
	url := s.slackSettings.ApiBaseURL + POST_MESSAGE_ENDPOINT
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	s.logger.Debug(
		"Sending POST to Slack",
		zap.String("url", url),
		zap.ByteString("body", jsonBody),
		zap.String("method", "POST"),
		zap.String("token", s.slackSettings.BotToken),
	)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.slackSettings.BotToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to post message: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)

	var slackResp models.ExternalPostMessageResponse
	if err := json.Unmarshal(bodyBytes, &slackResp); err != nil {
		s.logger.Error(
			"Failed to decode Slack response",
			zap.String("url", url),
			zap.String("raw_response", bodyString),
			zap.Error(err),
		)
		return nil, err
	}

	if !slackResp.Ok {
		s.logger.Error("Slack API returned ok=false", zap.String("url", url), zap.String("raw_response", bodyString))
		return nil, fmt.Errorf("slack api returned ok=false when posting message")
	}

	return &slackResp, nil
}

func (s *SlackClient) JoinConversation(
	ctx context.Context,
	request *models.ExternalJoinConversationRequest,
) (*models.ExternalJoinConversationResponse, error) {
	url := s.slackSettings.ApiBaseURL + JOIN_CONVERSATION_ENDPOINT
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	s.logger.Debug(
		"Sending POST to Slack",
		zap.String("url", url),
		zap.ByteString("body", jsonBody),
		zap.String("method", "POST"),
		zap.String("token", s.slackSettings.BotToken),
	)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.slackSettings.BotToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to join conversation: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)

	var slackResp models.ExternalJoinConversationResponse
	if err := json.Unmarshal(bodyBytes, &slackResp); err != nil {
		s.logger.Error(
			"Failed to decode Slack response",
			zap.String("url", url),
			zap.String("raw_response", bodyString),
			zap.Error(err),
		)
		return nil, err
	}

	if !slackResp.Ok {
		s.logger.Error("Slack API returned ok=false", zap.String("url", url), zap.String("raw_response", bodyString))
		return nil, fmt.Errorf("slack api returned ok=false when joining conversation")
	}

	return &slackResp, nil
}

func (s *SlackClient) GetConversationInfo(
	ctx context.Context,
	request *models.ExternalGetConversationInfoRequest,
) (*models.ExternalGetConversationInfoResponse, error) {
	url := s.slackSettings.ApiBaseURL + GET_CONVERSATION_INFO_ENDPOINT + "?channel=" + request.Channel

	s.logger.Debug(
		"Sending GET to Slack",
		zap.String("url", url),
		zap.String("method", "GET"),
		zap.String("token", s.slackSettings.BotToken),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.slackSettings.BotToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get conversation info: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)

	var slackResp models.ExternalGetConversationInfoResponse
	if err := json.Unmarshal(bodyBytes, &slackResp); err != nil {
		s.logger.Error(
			"Failed to decode Slack response",
			zap.String("url", url),
			zap.String("raw_response", bodyString),
			zap.Error(err),
		)
		return nil, err
	}

	if !slackResp.Ok {
		s.logger.Error("Slack API returned ok=false", zap.String("url", url), zap.String("raw_response", bodyString))
		return nil, fmt.Errorf("slack api returned ok=false when getting conversation info")
	}

	return &slackResp, nil
}

func (s *SlackClient) GetConversationMembers(
	ctx context.Context,
	request *models.ExternalGetConversationMembersRequest,
) (*models.ExternalGetConversationMembersResponse, error) {
	url := s.slackSettings.ApiBaseURL + GET_CONVERSATION_MEMBERS_ENDPOINT + "?channel=" + request.Channel

	s.logger.Debug(
		"Sending GET to Slack",
		zap.String("url", url),
		zap.String("method", "GET"),
		zap.String("token", s.slackSettings.BotToken),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.slackSettings.BotToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get conversation members: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)

	var slackResp models.ExternalGetConversationMembersResponse
	if err := json.Unmarshal(bodyBytes, &slackResp); err != nil {
		s.logger.Error(
			"Failed to decode Slack response",
			zap.String("url", url),
			zap.String("raw_response", bodyString),
			zap.Error(err),
		)
		return nil, err
	}

	if !slackResp.Ok {
		s.logger.Error("Slack API returned ok=false", zap.String("url", url), zap.String("raw_response", bodyString))
		return nil, fmt.Errorf("slack api returned ok=false when getting conversation members")
	}

	return &slackResp, nil
}

func (s *SlackClient) GetUsersInfo(
	ctx context.Context,
	request *models.ExternalGetUsersInfoRequest,
) (*models.ExternalGetUsersInfoResponse, error) {
	url := s.slackSettings.ApiBaseURL + GET_USERS_INFO_ENDPOINT + "?user=" + request.User

	s.logger.Debug(
		"Sending GET to Slack",
		zap.String("url", url),
		zap.String("method", "GET"),
		zap.String("token", s.slackSettings.BotToken),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.slackSettings.BotToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)

	var slackResp models.ExternalGetUsersInfoResponse
	if err := json.Unmarshal(bodyBytes, &slackResp); err != nil {
		s.logger.Error(
			"Failed to decode Slack response",
			zap.String("url", url),
			zap.String("raw_response", bodyString),
			zap.Error(err),
		)
		return nil, err
	}

	if !slackResp.Ok {
		s.logger.Error("Slack API returned ok=false", zap.String("url", url), zap.String("raw_response", bodyString))
		return nil, fmt.Errorf("slack api returned ok=false when getting user info")
	}

	return &slackResp, nil
}
