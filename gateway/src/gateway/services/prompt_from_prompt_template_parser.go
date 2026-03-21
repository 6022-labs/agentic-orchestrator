package services

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/google/uuid"
)

type PromptFromPromptTemplateParser struct {
	agentPromptTemplateRequester interfaces.AgentPromptTemplateRequesterInterface
}

func NewPromptFromPromptTemplateParser(
	agentPromptTemplateRequester interfaces.AgentPromptTemplateRequesterInterface,
) *PromptFromPromptTemplateParser {
	return &PromptFromPromptTemplateParser{
		agentPromptTemplateRequester: agentPromptTemplateRequester,
	}
}

func (p *PromptFromPromptTemplateParser) ParseCloneSystem(agentId uuid.UUID, parameters map[string]any) (*string, error) {
	promptTemplate, err := p.agentPromptTemplateRequester.GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyCloneSystem)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent prompt template: %w", err)
	}

	return p.parse(*promptTemplate, parameters)
}

func (p *PromptFromPromptTemplateParser) ParseCloneDefault(agentId uuid.UUID, parameters map[string]any) (*string, error) {
	promptTemplate, err := p.agentPromptTemplateRequester.GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyCloneDefault)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent prompt template: %w", err)
	}

	return p.parse(*promptTemplate, parameters)
}

func (p *PromptFromPromptTemplateParser) ParseFacilitatorSystem(agentId uuid.UUID, parameters map[string]any) (*string, error) {
	promptTemplate, err := p.agentPromptTemplateRequester.GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyFacilitatorSystem)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent prompt template: %w", err)
	}

	return p.parse(*promptTemplate, parameters)
}

func (p *PromptFromPromptTemplateParser) ParseFacilitatorDefault(agentId uuid.UUID, parameters map[string]any) (*string, error) {
	promptTemplate, err := p.agentPromptTemplateRequester.GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyFacilitatorDefault)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent prompt template: %w", err)
	}

	return p.parse(*promptTemplate, parameters)
}

func (p *PromptFromPromptTemplateParser) ParseFacilitatorConclude(agentId uuid.UUID, parameters map[string]any) (*string, error) {
	promptTemplate, err := p.agentPromptTemplateRequester.GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyFacilitatorConclude)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent prompt template: %w", err)
	}

	return p.parse(*promptTemplate, parameters)
}

func (p *PromptFromPromptTemplateParser) ParseExpertSystem(agentId uuid.UUID, parameters map[string]any) (*string, error) {
	promptTemplate, err := p.agentPromptTemplateRequester.GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyExpertSystem)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent prompt template: %w", err)
	}

	return p.parse(*promptTemplate, parameters)
}

func (p *PromptFromPromptTemplateParser) ParseExpertContextFound(agentId uuid.UUID, parameters map[string]any) (*string, error) {
	promptTemplate, err := p.agentPromptTemplateRequester.GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyExpertContextFound)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent prompt template: %w", err)
	}

	return p.parse(*promptTemplate, parameters)
}

func (p *PromptFromPromptTemplateParser) ParseExpertContextNotFound(agentId uuid.UUID, parameters map[string]any) (*string, error) {
	promptTemplate, err := p.agentPromptTemplateRequester.GetAgentPromptTemplateByAgentIdAndKey(agentId, enums.AgentPromptTemplateKeyExpertContextNotFound)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent prompt template: %w", err)
	}

	return p.parse(*promptTemplate, parameters)
}

func (p *PromptFromPromptTemplateParser) parse(promptTemplate string, parameters any) (*string, error) {
	promptTemplate = p.sanitizeTemplate(promptTemplate)

	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return nil, err
	}

	var outputBuffer strings.Builder
	err = tmpl.Execute(&outputBuffer, parameters)
	if err != nil {
		return nil, err
	}

	output := outputBuffer.String()
	return &output, nil
}

func (p *PromptFromPromptTemplateParser) sanitizeTemplate(promptTemplate string) string {
	var b strings.Builder
	for _, r := range promptTemplate {
		// Skip control chars and unprintable Unicode
		if r == '\u2028' || r == '\u2029' || r == '\u2060' || !unicode.IsPrint(r) {
			continue
		}

		if unicode.IsControl(r) {
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}
