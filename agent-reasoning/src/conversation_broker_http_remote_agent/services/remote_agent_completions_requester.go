package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_remote_agent/clients"
	http_node_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_remote_agent/models"
)

type RemoteAgentCompletionsRequester struct {
	client clients.RemoteAgentClientInterface
}

func NewRemoteAgentCompletionsRequester(client clients.RemoteAgentClientInterface) *RemoteAgentCompletionsRequester {
	return &RemoteAgentCompletionsRequester{client: client}
}

func (r *RemoteAgentCompletionsRequester) RequestCompletions(card *models.RemoteAgentCard, input *models.RemoteAgentInput) ([]string, error) {
	ensDomain := input.Metadata["ensDomain"]

	if card.Endpoints.Responses != "" {
		return r.requestResponsesTurn(card.Endpoints.Responses, ensDomain, input)
	}
	if card.Endpoints.ChatCompletions != "" {
		return r.requestChatCompletionTurn(card.Endpoints.ChatCompletions, ensDomain, input)
	}
	return nil, fmt.Errorf("card advertises no endpoint")
}

func (r *RemoteAgentCompletionsRequester) requestResponsesTurn(endpointUrl, ensDomain string, turn *models.RemoteAgentInput) ([]string, error) {
	input := make([]*http_node_models.ExternalResponseInputMessage, len(turn.Input))
	for i, message := range turn.Input {
		input[i] = &http_node_models.ExternalResponseInputMessage{
			Type:    "message",
			Role:    message.Role.String(),
			Content: message.Content,
		}
	}

	response, err := r.client.CreateResponse(endpointUrl, ensDomain, &http_node_models.ExternalCreateResponseRequest{
		Instructions: turn.Instructions,
		Input:        input,
		Metadata:     turn.Metadata,
	})
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("agent node returned error %s: %s", response.Error.Code, response.Error.Message)
	}

	completions := make([]string, 0, len(response.Output))
	for _, output := range response.Output {
		if output.Type != "message" {
			continue
		}
		for _, content := range output.Content {
			if content.Text != "" {
				completions = append(completions, content.Text)
			}
		}
	}

	if len(completions) == 0 {
		return nil, fmt.Errorf("agent node returned no message output (status %s)", response.Status)
	}

	return completions, nil
}

func (r *RemoteAgentCompletionsRequester) requestChatCompletionTurn(endpointUrl, ensDomain string, turn *models.RemoteAgentInput) ([]string, error) {
	messages := make([]*http_node_models.ExternalChatCompletionMessage, 0, len(turn.Input)+1)
	if turn.Instructions != "" {
		messages = append(messages, &http_node_models.ExternalChatCompletionMessage{
			Role:    "system",
			Content: turn.Instructions,
		})
	}
	for _, message := range turn.Input {
		messages = append(messages, &http_node_models.ExternalChatCompletionMessage{
			Role:    message.Role.String(),
			Content: message.Content,
		})
	}

	response, err := r.client.CreateChatCompletion(endpointUrl, ensDomain, &http_node_models.ExternalChatCompletionRequest{
		Messages: messages,
		Metadata: turn.Metadata,
	})
	if err != nil {
		return nil, err
	}

	completions := make([]string, 0, len(response.Choices))
	for _, choice := range response.Choices {
		if choice.Message != nil && choice.Message.Content != "" {
			completions = append(completions, choice.Message.Content)
		}
	}

	if len(completions) == 0 {
		return nil, fmt.Errorf("agent node returned no chat completion choice")
	}

	return completions, nil
}
