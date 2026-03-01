package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
)

type ChatCompletionRequesterInterface interface {
	ChatCompletion(
		llmConnection *models.LLMConnection,
		mcpConnections []*models.MCPConnection,
		messages []*models.ChatCompletionMessage,
		apiKey *string,
	) (*string, error)
}
