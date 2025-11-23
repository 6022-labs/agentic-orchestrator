package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type CreateAgentSecondStepInterface interface {
	Execute(request *requests.AgentByIdRequest) (*responses.CreateAgentSecondStepResponse, error)
}

type CreateAgentSecondStep struct {
	orchestratorSettings *settings.OrchestratorSettings
	mintSigner           interfaces.MintSignerInterface
	agentRepository      repositories.AgentRepositoryInterface
}

func NewCreateAgentSecondStep(
	mintSigner interfaces.MintSignerInterface,
	orchestratorSettings *settings.OrchestratorSettings,
	agentRepository repositories.AgentRepositoryInterface,
) *CreateAgentSecondStep {
	return &CreateAgentSecondStep{
		mintSigner:           mintSigner,
		agentRepository:      agentRepository,
		orchestratorSettings: orchestratorSettings,
	}
}

func (f *CreateAgentSecondStep) Execute(request *requests.AgentByIdRequest) (*responses.CreateAgentSecondStepResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	// Fetch the agent from the database
	agent, err := f.agentRepository.GetById(request.AgentId)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, errors.NewNotFoundError("Agent not found")
	}
	if agent.CollectionTokenId != nil {
		return nil, errors.NewValidationError("agentId", "Agent is already minted")
	}
	if agent.Creator != request.AuthenticatedWallet.Hex() {
		return nil, errors.NewValidationError("agentCollectionTokenId", "You are not the creator of this agent")
	}

	generatedNonce := uuid.New()
	agentWallets := []string{agent.WalletAddress}

	signature, err := f.mintSigner.Sign(
		request.AuthenticatedWallet.Hex(),
		agent.Name,
		agent.Role.String(),
		agent.CloneOf,
		agentWallets,
		agent.Images,
		agent.Attributes,
		generatedNonce,
		common.HexToAddress(agent.CollectionAddress),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign agent creation: %w", err)
	}

	return &responses.CreateAgentSecondStepResponse{
		Signature: *signature,
		Nonce:     generatedNonce.String(),
		Signer:    f.orchestratorSettings.PublicKey.Hex(),
	}, nil
}
