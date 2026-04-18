package model_builders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type ParticipantBuilder struct {
	id                uuid.UUID
	collectionAddress *string
	collectionTokenId *uint64
	ensDomain         *string
	role              *enums.ParticipantRole
	displayName       *string
	isProcessing      bool
	lastProcessedAt   *time.Time
	createdAt         time.Time
	updatedAt         time.Time
}

// NewAgentParticipantBuilder builds an agent participant with role=clone by default.
func NewAgentParticipantBuilder() *ParticipantBuilder {
	collectionAddress := fmt.Sprintf("0x%040x", rand.Uint64())
	collectionTokenId := rand.Uint64()
	ensDomain := faker.Word() + ".eth"
	role := enums.ParticipantRoleClone

	return &ParticipantBuilder{
		id:                uuid.New(),
		collectionAddress: &collectionAddress,
		collectionTokenId: &collectionTokenId,
		ensDomain:         &ensDomain,
		role:              &role,
		createdAt:         time.Now(),
		updatedAt:         time.Now(),
	}
}

// NewUserParticipantBuilder builds a user participant (display-name only).
func NewUserParticipantBuilder() *ParticipantBuilder {
	displayName := faker.Name()
	return &ParticipantBuilder{
		id:          uuid.New(),
		displayName: &displayName,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

func (b *ParticipantBuilder) WithID(id uuid.UUID) *ParticipantBuilder {
	b.id = id
	return b
}

func (b *ParticipantBuilder) WithRole(role enums.ParticipantRole) *ParticipantBuilder {
	b.role = &role
	return b
}

func (b *ParticipantBuilder) WithIsProcessing(v bool) *ParticipantBuilder {
	b.isProcessing = v
	return b
}

func (b *ParticipantBuilder) WithLastProcessedAt(t time.Time) *ParticipantBuilder {
	b.lastProcessedAt = &t
	return b
}

func (b *ParticipantBuilder) WithDisplayName(name string) *ParticipantBuilder {
	b.displayName = &name
	return b
}

func (b *ParticipantBuilder) Build() *models.Participant {
	return &models.Participant{
		ID:                b.id,
		CollectionAddress: b.collectionAddress,
		CollectionTokenId: b.collectionTokenId,
		EnsDomain:         b.ensDomain,
		Role:              b.role,
		DisplayName:       b.displayName,
		IsProcessing:      b.isProcessing,
		LastProcessedAt:   b.lastProcessedAt,
		CreatedAt:         b.createdAt,
		UpdatedAt:         b.updatedAt,
	}
}
