package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/compute"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/utils"
	"go.uber.org/zap"
)

type AgentComputeInterface interface {
	Compute(request *requests.AgentComputeRequest) (*responses.AgentComputeResponses, error)
}

type AgentCompute struct {
	logger                     *zap.Logger
	agentRegistryRequester     interfaces.AgentRegistryRequesterInterface
	cloneComputeStrategy       compute.CloneComputeStrategyInterface
	humanComputeStrategy       compute.HumanComputeStrategyInterface
	expertComputeStrategy      compute.ExpertComputeStrategyInterface
	facilitatorComputeStrategy compute.FacilitatorComputeStrategyInterface
}

func NewAgentCompute(
	logger *zap.Logger,
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	cloneComputeStrategy compute.CloneComputeStrategyInterface,
	humanComputeStrategy compute.HumanComputeStrategyInterface,
	expertComputeStrategy compute.ExpertComputeStrategyInterface,
	facilitatorComputeStrategy compute.FacilitatorComputeStrategyInterface,
) *AgentCompute {
	return &AgentCompute{
		logger:                     logger,
		agentRegistryRequester:     agentRegistryRequester,
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

	agentResponse, err := a.agentRegistryRequester.GetAgent(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}
	if agentResponse == nil {
		return nil, fmt.Errorf("agent with collection token ID %d not found", request.AgentCollectionTokenId)
	}
	agent := agentResponseToModel(agentResponse)

	isFinal := false

	completions, final, err := a.computeByRole(agent, request)
	if err != nil {
		return nil, fmt.Errorf("failed to compute agent %s: %w", agent.EnsDomain, err)
	}

	if final != nil {
		isFinal = *final
	}

	return &responses.AgentComputeResponses{
		IsFinal:   isFinal,
		Responses: completions,
	}, nil
}

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
			Content:        utils.DemojizeContent(input.Content),
			CreatedAt:      input.Timestamp,
			Platform:       request.Platform,
			ConversationId: &request.ConversationContext.Id,
		}
	}

	return agentContextMessages
}

func agentResponseToModel(r *responses.AgentResponse) *models.Agent {
	return &models.Agent{
		ID:                r.ID,
		Name:              r.Name,
		Role:              r.Role,
		Owner:             r.Owner,
		Creator:           r.Creator,
		Images:            r.Images,
		CloneOf:           r.CloneOf,
		Attributes:        r.Attributes,
		EnsDomain:         r.EnsDomain,
		IsInternal:        r.IsInternal,
		WalletAddress:     r.WalletAddress,
		ProposalId:        r.ProposalId,
		CollectionTokenId: r.CollectionTokenId,
		CollectionAddress: r.CollectionAddress,
	}
}
