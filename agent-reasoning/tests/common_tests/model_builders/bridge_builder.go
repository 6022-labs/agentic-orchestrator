package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type BridgeBuilder struct {
	id          uuid.UUID
	kind        enums.BridgeKind
	callbackUrl string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewBridgeBuilder() *BridgeBuilder {
	return &BridgeBuilder{
		id:          uuid.New(),
		kind:        enums.BridgeKindSlack,
		callbackUrl: faker.URL(),
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

func (b *BridgeBuilder) WithID(id uuid.UUID) *BridgeBuilder {
	b.id = id
	return b
}

func (b *BridgeBuilder) WithKind(kind enums.BridgeKind) *BridgeBuilder {
	b.kind = kind
	return b
}

func (b *BridgeBuilder) WithCallbackUrl(url string) *BridgeBuilder {
	b.callbackUrl = url
	return b
}

func (b *BridgeBuilder) Build() *models.Bridge {
	return &models.Bridge{
		ID:          b.id,
		Kind:        b.kind,
		CallbackUrl: b.callbackUrl,
		CreatedAt:   b.createdAt,
		UpdatedAt:   b.updatedAt,
	}
}
