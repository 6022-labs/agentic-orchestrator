package model_builders

import (
	"strings"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type AgentAddressBuilder struct {
	agentId     uuid.UUID
	addressType enums.AddressType
	address     string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewAgentAddressBuilder() *AgentAddressBuilder {
	createdAtStr := faker.Date()
	updatedAtStr := faker.Date()

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
	updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

	return &AgentAddressBuilder{
		agentId:     uuid.New(),
		addressType: enums.AddressTypeEvm,
		address:     strings.ToLower(common_tests.RandomEthAddress().Hex()),
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (b *AgentAddressBuilder) WithAgentId(agentId uuid.UUID) *AgentAddressBuilder {
	b.agentId = agentId
	return b
}

func (b *AgentAddressBuilder) WithType(addressType enums.AddressType) *AgentAddressBuilder {
	b.addressType = addressType
	return b
}

func (b *AgentAddressBuilder) WithAddress(address string) *AgentAddressBuilder {
	b.address = address
	return b
}

func (b *AgentAddressBuilder) Build() *models.AgentAddress {
	return &models.AgentAddress{
		AgentId:   b.agentId,
		Type:      b.addressType,
		Address:   b.address,
		CreatedAt: b.createdAt,
		UpdatedAt: b.updatedAt,
	}
}
