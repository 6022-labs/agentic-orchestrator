package services

import (
	"encoding/json"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/clients"
	http_slack_models "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/settings"
)

type SlackThreadRequester struct {
	client              clients.SlackClientInterface
	ipfsGatewaySettings *settings.IpfsGatewaySettings
}

func NewSlackThreadRequester(
	client clients.SlackClientInterface,
	ipfsGatewaySettings *settings.IpfsGatewaySettings,
) *SlackThreadRequester {
	return &SlackThreadRequester{
		client:              client,
		ipfsGatewaySettings: ipfsGatewaySettings,
	}
}

func (s *SlackThreadRequester) PushCompletionToThread(completion string, slackAgent *models.SlackAgent, slackChannelId, slackThreadId string) error {
	request := &http_slack_models.ExternalPostMessageRequest{
		Channel:  slackChannelId,
		Text:     completion,
		ThreadTs: &slackThreadId,
	}

	if slackAgent != nil {
		request.Username = &slackAgent.EnsDomain

		metadataObj := map[string]interface{}{
			"event_type": event_handlers.AGENT_MESSAGE_CREATED_METADATA_EVENT_TYPE,
			"event_payload": map[string]interface{}{
				"collection_token_address": slackAgent.CollectionAddress,
				"collection_token_id":      slackAgent.CollectionTokenId,
			},
		}

		metadataBytes, err := json.Marshal(metadataObj)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}

		metadataStr := string(metadataBytes)
		request.Metadata = &metadataStr
	}

	if slackAgent != nil && slackAgent.IconImageCid != nil {
		iconUrl := fmt.Sprintf(s.ipfsGatewaySettings.GatewayUrlTemplate, *slackAgent.IconImageCid)
		request.IconUrl = &iconUrl
	}

	_, err := s.client.PostMessage(request)
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}
