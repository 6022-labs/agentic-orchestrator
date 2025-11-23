package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/compute"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AgentComputeInterface interface {
	Compute(request *requests.AgentComputeRequest) (*responses.AgentComputeResponses, error)
}

type AgentCompute struct {
	logger                     *zap.Logger
	orchestratorSettings       *settings.OrchestratorSettings
	agentRepository            repositories.AgentRepositoryInterface
	cloneComputeStrategy       compute.CloneComputeStrategyInterface
	humanComputeStrategy       compute.HumanComputeStrategyInterface
	expertComputeStrategy      compute.ExpertComputeStrategyInterface
	facilitatorComputeStrategy compute.FacilitatorComputeStrategyInterface
}

func NewAgentCompute(
	logger *zap.Logger,
	orchestratorSettings *settings.OrchestratorSettings,
	agentRepository repositories.AgentRepositoryInterface,
	cloneComputeStrategy compute.CloneComputeStrategyInterface,
	humanComputeStrategy compute.HumanComputeStrategyInterface,
	expertComputeStrategy compute.ExpertComputeStrategyInterface,
	facilitatorComputeStrategy compute.FacilitatorComputeStrategyInterface,
) *AgentCompute {
	return &AgentCompute{
		logger:                     logger,
		agentRepository:            agentRepository,
		orchestratorSettings:       orchestratorSettings,
		cloneComputeStrategy:       cloneComputeStrategy,
		humanComputeStrategy:       humanComputeStrategy,
		expertComputeStrategy:      expertComputeStrategy,
		facilitatorComputeStrategy: facilitatorComputeStrategy,
	}
}

func (a *AgentCompute) Compute(request *requests.AgentComputeRequest) (*responses.AgentComputeResponses, error) {
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, err
	}

	if request.AuthenticatedWallet.Hex() != a.orchestratorSettings.PublicKey.Hex() {
		return nil, fmt.Errorf("only the orchestrator himself can perform agent computations")
	}

	agent, err := a.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if agent == nil {
		return nil, fmt.Errorf("agent with collection token ID %d not found", request.AgentCollectionTokenId)
	}

	isFinal := false

	completions, final, err := a.computeByRole(agent, request)
	if err != nil {
		return nil, fmt.Errorf("failed to compute agent %s: %w", agent.EnsDomain, err)
	}

	// If facilitator, set isFinal
	if final != nil {
		isFinal = *final
	}

	return &responses.AgentComputeResponses{
		IsFinal:   isFinal,
		Responses: completions,
	}, nil
}

// computeByRole selects the compute strategy based on agent role.
func (a *AgentCompute) computeByRole(agent *models.Agent, request *requests.AgentComputeRequest) ([]string, *bool, error) {
	agentContextMessages := a.requestToAgentContextMessages(agent, request)
	context := request.ToComputeContext()

	switch agent.Role {
	case enums.AgentRoleHuman:
		completions, err := a.humanComputeStrategy.Compute(agent, agentContextMessages, context)
		return completions, nil, err
	case enums.AgentRoleClone:
		completions, err := a.cloneComputeStrategy.Compute(agent, agentContextMessages, context)
		return completions, nil, err
	case enums.AgentRoleFacilitator:
		completions, isFinal, err := a.facilitatorComputeStrategy.Compute(agent, agentContextMessages, context)
		return completions, isFinal, err
	default:
		completions, err := a.expertComputeStrategy.Compute(agent, agentContextMessages, context)
		return completions, nil, err
	}
}

func (a *AgentCompute) requestToAgentContextMessages(
	agent *models.Agent,
	request *requests.AgentComputeRequest,
) []*models.AgentContextMessage {
	agentContextMessages := make([]*models.AgentContextMessage, len(request.Inputs))

	for index, input := range request.Inputs {
		role := enums.AgentContextMessageRoleUser
		if agent.Role != enums.AgentRoleHuman && input.Author == agent.EnsDomain {
			role = enums.AgentContextMessageRoleAssistant
		}

		agentContextMessages[index] = &models.AgentContextMessage{
			Role:           role,
			Author:         input.Author,
			Content:        input.Content,
			CreatedAt:      input.Timestamp,
			Platform:       request.Platform,
			ConversationId: &request.ConversationContext.Id,
		}
	}

	return agentContextMessages
}
