package models

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/interfaces"

type ExternalMcpToolsResponse struct {
	Tools []ExternalMcpToolsResponseTool `json:"tools"`
}

type ExternalMcpToolsResponseTool struct {
	Name        string                              `json:"name"`
	Title       string                              `json:"title"`
	Description string                              `json:"description"`
	InputSchema ExternalMcpToolsResponseInputSchema `json:"inputSchema"`
}

type ExternalMcpToolsResponseInputSchema struct {
	Type                 string                                                 `json:"type"`
	Properties           map[string]ExternalMcpToolsResponseInputSchemaProperty `json:"properties"`
	Required             []string                                               `json:"required"`
	AdditionalProperties bool                                                   `json:"additionalProperties"`
	Schema               string                                                 `json:"$schema"`
}

type ExternalMcpToolsResponseInputSchemaProperty struct {
	Type        string                                       `json:"type"`
	Items       *ExternalMcpToolsResponseInputSchemaProperty `json:"items,omitempty"`
	Description *string                                      `json:"description,omitempty"`
}

func (e *ExternalMcpToolsResponse) ToMcpTools() []*interfaces.McpTool {
	tools := make([]*interfaces.McpTool, len(e.Tools))
	for index, tool := range e.Tools {
		tools[index] = &interfaces.McpTool{
			Name:        tool.Name,
			Description: tool.Description,
			InputScheme: interfaces.McpToolInputScheme{
				Type:                 tool.InputSchema.Type,
				Required:             tool.InputSchema.Required,
				AdditionalProperties: tool.InputSchema.AdditionalProperties,
				Schema:               tool.InputSchema.Schema,
			},
		}

		for propName, prop := range tool.InputSchema.Properties {
			if tools[index].InputScheme.Properties == nil {
				tools[index].InputScheme.Properties = make(map[string]interfaces.ToolInputSchemeProperty)
			}

			tools[index].InputScheme.Properties[propName] = prop.ToInternalModel()
		}
	}

	return tools
}

func (p ExternalMcpToolsResponseInputSchemaProperty) ToInternalModel() interfaces.ToolInputSchemeProperty {
	this := interfaces.ToolInputSchemeProperty{
		Type:        p.Type,
		Description: p.Description,
	}

	if p.Items != nil {
		internalItems := p.Items.ToInternalModel()
		this.Items = &internalItems
	}

	return this
}
