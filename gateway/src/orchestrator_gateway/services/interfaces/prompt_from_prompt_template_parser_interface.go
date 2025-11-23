package interfaces

import "github.com/google/uuid"

type PromptFromPromptTemplateParserInterface interface {
	ParseCloneSystem(agentId uuid.UUID, parameters map[string]any) (*string, error)
	ParseCloneDefault(agentId uuid.UUID, parameters map[string]any) (*string, error)
	ParseFacilitatorSystem(agentId uuid.UUID, parameters map[string]any) (*string, error)
	ParseFacilitatorDefault(agentId uuid.UUID, parameters map[string]any) (*string, error)
	ParseFacilitatorConclude(agentId uuid.UUID, parameters map[string]any) (*string, error)
	ParseExpertSystem(agentId uuid.UUID, parameters map[string]any) (*string, error)
	ParseExpertContextFound(agentId uuid.UUID, parameters map[string]any) (*string, error)
	ParseExpertContextNotFound(agentId uuid.UUID, parameters map[string]any) (*string, error)
}
