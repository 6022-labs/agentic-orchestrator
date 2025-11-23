package models

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/services/interfaces"

type ExternalToolInputSchemaProperty struct {
	Type        string                           `json:"type"`
	Items       *ExternalToolInputSchemaProperty `json:"items,omitempty"`
	Description *string                          `json:"description,omitempty"`
}

type ExternalToolInputSchemaProperties map[string]ExternalToolInputSchemaProperty

func NewExternalToolInputSchemaPropertyFromInternalModel(internalProperty interfaces.ToolInputSchemeProperty) ExternalToolInputSchemaProperty {
	this := ExternalToolInputSchemaProperty{
		Type:        internalProperty.Type,
		Description: internalProperty.Description,
	}

	if internalProperty.Items != nil {
		internalItems := NewExternalToolInputSchemaPropertyFromInternalModel(*internalProperty.Items)
		this.Items = &internalItems
	}

	return this
}

func NewExternalToolInputSchemaPropertiesFromInternalModels(internalProperties map[string]interfaces.ToolInputSchemeProperty) ExternalToolInputSchemaProperties {
	externalProperties := make(ExternalToolInputSchemaProperties)

	for key, internalProperty := range internalProperties {
		externalProperties[key] = NewExternalToolInputSchemaPropertyFromInternalModel(internalProperty)
	}

	return externalProperties
}

func (p *ExternalToolInputSchemaProperty) ToInternalModel() interfaces.ToolInputSchemeProperty {
	internalModel := interfaces.ToolInputSchemeProperty{
		Type:        p.Type,
		Description: p.Description,
	}

	if p.Items != nil {
		internalItems := p.Items.ToInternalModel()
		internalModel.Items = &internalItems
	}

	return internalModel
}

func (p ExternalToolInputSchemaProperties) ToInternalModels() map[string]interfaces.ToolInputSchemeProperty {
	internalProperties := make(map[string]interfaces.ToolInputSchemeProperty)

	for key, externalProperty := range p {
		internalProperties[key] = externalProperty.ToInternalModel()
	}

	return internalProperties
}
