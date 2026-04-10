package use_cases

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"gorm.io/gorm"
)

type LLMConnectionChatCompletionInterface interface {
	Complete(request *requests.LLMConnectionChatCompletionRequest) (*string, error)
}

type LLMConnectionChatCompletion struct {
	agentLLMSecretStorage               services.AgentLLMSecretStorageInterface
	chatCompletionRequester             interfaces.ChatCompletionRequesterInterface
	llmConnectionRepository             repositories.LLMConnectionRepositoryInterface
	mcpConnectionRepository             repositories.MCPConnectionRepositoryInterface
	mcpConnectionHeaderAuthValueStorage services.MCPConnectionHeaderAuthValueStorageInterface
}

func NewLLMConnectionChatCompletion(
	agentLLMSecretStorage services.AgentLLMSecretStorageInterface,
	chatCompletionRequester interfaces.ChatCompletionRequesterInterface,
	llmConnectionRepository repositories.LLMConnectionRepositoryInterface,
	mcpConnectionRepository repositories.MCPConnectionRepositoryInterface,
	mcpConnectionHeaderAuthValueStorage services.MCPConnectionHeaderAuthValueStorageInterface,
) *LLMConnectionChatCompletion {
	return &LLMConnectionChatCompletion{
		agentLLMSecretStorage:               agentLLMSecretStorage,
		chatCompletionRequester:             chatCompletionRequester,
		llmConnectionRepository:             llmConnectionRepository,
		mcpConnectionRepository:             mcpConnectionRepository,
		mcpConnectionHeaderAuthValueStorage: mcpConnectionHeaderAuthValueStorage,
	}
}

func (uc *LLMConnectionChatCompletion) Complete(request *requests.LLMConnectionChatCompletionRequest) (*string, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	llmConnection, err := uc.llmConnectionRepository.GetByAgentId(request.AgentId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get llm connection: %w", err)
	}
	if llmConnection == nil {
		return nil, fmt.Errorf("llm connection not found for agentId: %s", request.AgentId)
	}

	mcpConnections, err := uc.mcpConnectionRepository.GetByAgentId(request.AgentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get mcp connections: %w", err)
	}

	for _, mcpConnection := range mcpConnections {
		headerAuth, err := uc.mcpConnectionHeaderAuthValueStorage.Get(mcpConnection)
		if err != nil {
			return nil, fmt.Errorf("failed to get header auth for agentId %s for MCPConnection: %w", request.AgentId, err)
		}

		mcpConnection.HeaderAuth = headerAuth
	}

	var apiKey *string
	if llmConnection.Platform != enums.LLMConnectionPlatformOllama {
		apiKey, err = uc.agentLLMSecretStorage.Get(request.AgentId)
		if err != nil {
			return nil, fmt.Errorf("failed to get api key for agentId %s: %w", request.AgentId, err)
		}
	}

	chatCompletion, err := uc.chatCompletionRequester.ChatCompletion(
		llmConnection,
		mcpConnections,
		request.Messages,
		apiKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat completion for agentId %s: %w", request.AgentId, err)
	}

	return chatCompletion, err
}
