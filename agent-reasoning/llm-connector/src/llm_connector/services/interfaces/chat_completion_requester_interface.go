package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models"
)

type ChatCompletionRequesterInterface interface {
	ChatCompletion(
		llmConnection *models.LLMConnection,
		mcpConnections []*models.MCPConnection,
		messages []*models.ChatCompletionMessage,
		apiKey *string,
	) (*string, error)
}
