package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
)

// AbstractAgentAddressEventHandler carries typed-address mapping and lookup logic shared by synchronizer handlers.
type AbstractAgentAddressEventHandler struct {
	agentRepository        repositories.AgentRepositoryInterface
	agentAddressRepository repositories.AgentAddressRepositoryInterface
}

func NewAbstractAgentAddressEventHandler(
	agentRepository repositories.AgentRepositoryInterface,
	agentAddressRepository repositories.AgentAddressRepositoryInterface,
) *AbstractAgentAddressEventHandler {
	return &AbstractAgentAddressEventHandler{
		agentRepository:        agentRepository,
		agentAddressRepository: agentAddressRepository,
	}
}

// typedAddresses maps on-chain entries to one address per known type (first entry wins, unknown types skipped).
func (handler *AbstractAgentAddressEventHandler) typedAddresses(
	entries []abi.AgentAddress,
) map[enums.AddressType]string {
	addresses := map[enums.AddressType]string{}
	for _, entry := range entries {
		var addressType enums.AddressType
		if err := addressType.Set(entry.AddressType); err != nil {
			continue
		}
		if _, exists := addresses[addressType]; !exists {
			addresses[addressType] = entry.Value
		}
	}
	return addresses
}

// findAgentByAddresses returns the agent holding any of the typed addresses, or nil.
func (handler *AbstractAgentAddressEventHandler) findAgentByAddresses(
	addresses map[enums.AddressType]string,
) (*models.Agent, error) {
	for _, addressType := range enums.GetAllAddressTypes() {
		value, exists := addresses[addressType]
		if !exists {
			continue
		}
		agentAddress, err := handler.agentAddressRepository.GetByTypeAndAddress(addressType, value)
		if err != nil {
			return nil, fmt.Errorf("failed to get agent address: %w", err)
		}
		if agentAddress == nil {
			continue
		}

		agent, err := handler.agentRepository.GetById(agentAddress.AgentId)
		if err != nil {
			return nil, fmt.Errorf("failed to get agent by id: %w", err)
		}
		if agent == nil {
			continue
		}
		return agent, nil
	}

	return nil, nil
}
