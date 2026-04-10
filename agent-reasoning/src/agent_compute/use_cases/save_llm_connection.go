package use_cases

import (
	"errors"
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"gorm.io/gorm"
)

const (
	AMAZON_BEDROCK_HOST_TEMPLATE = "https://bedrock-runtime.%s.amazonaws.com"
)

type SaveLLMConnectionInterface interface {
	Save(request *requests.SaveLLMConnectionRequest) error
}

type SaveLLMConnection struct {
	agentLLMSecretStorage   services.AgentLLMSecretStorageInterface
	llmConnectionRepository repositories.LLMConnectionRepositoryInterface
}

func NewSaveLLMConnection(
	agentLLMSecretStorage services.AgentLLMSecretStorageInterface,
	llmConnectionRepository repositories.LLMConnectionRepositoryInterface,
) *SaveLLMConnection {
	return &SaveLLMConnection{
		agentLLMSecretStorage:   agentLLMSecretStorage,
		llmConnectionRepository: llmConnectionRepository,
	}
}

func (uc *SaveLLMConnection) Save(request *requests.SaveLLMConnectionRequest) error {
	// Validate the request
	if err := request.Validate(); err != nil {
		return err
	}

	connection, err := uc.llmConnectionRepository.GetByAgentId(request.AgentId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error requesting existing connection: %w", err)
	}
	if connection == nil {
		// Create the LLM connection
		connection = &models.LLMConnection{
			AgentID: request.AgentId,
		}
	}

	connection.Model = request.Model
	connection.Platform = request.Platform
	connection.SetupStatus = enums.LLMConnectionSetupStatusEstablished

	connection.BaseUrl = nil

	if request.BaseUrl != nil {
		sanitizedBaseUrl := strings.TrimSuffix(*request.BaseUrl, "/")
		connection.BaseUrl = &sanitizedBaseUrl
	}

	if request.Platform == enums.LLMConnectionPlatformAmazonBedrock {
		baseUrl := fmt.Sprintf(AMAZON_BEDROCK_HOST_TEMPLATE, *request.AmazonPlatformRegion)
		connection.BaseUrl = &baseUrl
	}

	_, err = uc.llmConnectionRepository.Save(connection)
	if err != nil {
		return fmt.Errorf("error saving LLM connection: %w", err)
	}

	if enums.GetAllCloudLLMConnectionPlatforms().Contains(request.Platform) {
		err = uc.agentLLMSecretStorage.Store(request.AgentId, *request.ApiKey)
		if err != nil {
			return fmt.Errorf("error storing API key: %w", err)
		}
	} else {
		uc.agentLLMSecretStorage.Delete(request.AgentId)
	}

	return nil
}
