package responses

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/models/enums"

type GetLLMConnectionSetupStatusResponse struct {
	IsCompleted bool                            `json:"isCompleted"`
	SetupStatus *enums.LLMConnectionSetupStatus `json:"setupStatus"`
}
