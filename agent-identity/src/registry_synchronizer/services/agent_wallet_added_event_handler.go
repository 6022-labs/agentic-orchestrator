package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"go.uber.org/zap"
)

type AgentWalletAddedEventHandler struct {
	logger                             *zap.Logger
	agentRepository                    repositories.AgentRepositoryInterface
	agentCollectionBlockchainRequester interfaces.AgentCollectionBlockchainRequesterInterface
}

func NewAgentWalletAddedEventHandler(
	logger *zap.Logger,
	agentRepository repositories.AgentRepositoryInterface,
	agentCollectionBlockchainRequester interfaces.AgentCollectionBlockchainRequesterInterface,
) *AgentWalletAddedEventHandler {
	return &AgentWalletAddedEventHandler{
		logger: logger,

		agentRepository:                    agentRepository,
		agentCollectionBlockchainRequester: agentCollectionBlockchainRequester,
	}
}

func (handler *AgentWalletAddedEventHandler) Handle(event *abi.AgentCollectionV1AgentWalletAdded) error {
	handler.logger.Info("Handling AgentWalletAdded event",
		zap.String("collectionAddress", event.Raw.Address.Hex()),
		zap.Uint64("tokenId", event.TokenId.Uint64()),
		zap.String("wallet", event.Wallet.Hex()),
	)

	agent, err := handler.agentRepository.GetByCollectionAddressAndCollectionTokenId(event.Raw.Address.Hex(), event.TokenId.Uint64())
	if err != nil {
		return fmt.Errorf("failed to get agent by collection address and token id: %w", err)
	}

	agentInformation, err := handler.agentCollectionBlockchainRequester.GetAgentInformation(event.Raw.Address, *event.TokenId)
	if err != nil {
		return fmt.Errorf("failed to get agent information: %w", err)
	}

	walletAddress := ""
	if len(agentInformation.Wallets) > 0 {
		walletAddress = agentInformation.Wallets[0].Hex()
	}

	if agent.WalletAddress == walletAddress {
		return nil
	}

	agent.WalletAddress = walletAddress
	if _, err := handler.agentRepository.Save(agent); err != nil {
		return fmt.Errorf("failed to save agent: %w", err)
	}

	return nil
}
