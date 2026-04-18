package model_builders

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type SlackChannelBuilder struct {
	id          string
	name        string
	description string
	agents      []models.SlackAgent
	users       []models.SlackUser
}

func NewSlackChannelBuilder() *SlackChannelBuilder {
	return &SlackChannelBuilder{
		id:          uuid.New().String(),
		name:        faker.Name(),
		description: faker.Sentence(),
		agents:      []models.SlackAgent{},
		users:       []models.SlackUser{},
	}
}

func (b *SlackChannelBuilder) WithId(id string) *SlackChannelBuilder {
	b.id = id
	return b
}

func (b *SlackChannelBuilder) WithName(name string) *SlackChannelBuilder {
	b.name = name
	return b
}

func (b *SlackChannelBuilder) WithDescription(description string) *SlackChannelBuilder {
	b.description = description
	return b
}

func (b *SlackChannelBuilder) WithAgents(agents []models.SlackAgent) *SlackChannelBuilder {
	b.agents = agents
	return b
}

func (b *SlackChannelBuilder) WithUsers(users []models.SlackUser) *SlackChannelBuilder {
	b.users = users
	return b
}

func (b *SlackChannelBuilder) Build() *models.SlackChannel {
	return &models.SlackChannel{
		Id:          b.id,
		Name:        b.name,
		Description: b.description,
		Agents:      b.agents,
		Users:       b.users,
	}
}
