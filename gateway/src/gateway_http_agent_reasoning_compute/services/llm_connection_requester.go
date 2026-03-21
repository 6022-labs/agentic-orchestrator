package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_reasoning_compute/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_reasoning_compute/models"
	"github.com/google/uuid"
)

type LLMConnectionRequester struct {
	client clients.AgentReasoningComputeClientInterface
}

func NewLLMConnectionRequester(
	client clients.AgentReasoningComputeClientInterface,
) *LLMConnectionRequester {
	return &LLMConnectionRequester{
		client: client,
	}
}

func (requester *LLMConnectionRequester) Save(request *interfaces.SaveLLMConnectionRequest) error {
	err := requester.client.SaveLLMConnection(&models.ExternalSaveLLMConnectionRequest{
		Model:    request.Model,
		ApiKey:   request.ApiKey,
		BaseUrl:  request.BaseUrl,
		AgentId:  request.AgentId,
		Platform: request.Platform,

		AmazonPlatformRegion: request.AmazonPlatformRegion,
	})
	if err != nil {
		return fmt.Errorf("failed to save LLM connection: %w", err)
	}

	return nil
}

func (requester *LLMConnectionRequester) Remove(agentId uuid.UUID) error {
	err := requester.client.RemoveLLMConnection(&models.ExternalAgentIdRequest{
		AgentId: agentId,
	})
	if err != nil {
		return fmt.Errorf("failed to remove LLM connection: %w", err)
	}

	return nil
}

func (requester *LLMConnectionRequester) GetSettings(agentId uuid.UUID) (*interfaces.LLMConnectionSettings, error) {
	response, err := requester.client.GetLLMConnectionSettings(&models.ExternalAgentIdRequest{
		AgentId: agentId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM connection settings: %w", err)
	}

	return response.ToLLMConnectionSettings(), nil
}

func (requester *LLMConnectionRequester) GetSetupStatus(agentId uuid.UUID) (*interfaces.LLMConnectionStatus, error) {
	response, err := requester.client.GetLLMConnectionSetupStatus(&models.ExternalAgentIdRequest{
		AgentId: agentId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM connection setup status: %w", err)
	}

	return &interfaces.LLMConnectionStatus{
		IsCompleted: response.IsCompleted,
		Status:      response.SetupStatus,
	}, nil
}
