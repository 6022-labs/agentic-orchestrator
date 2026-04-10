package responses

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"

type GetLLMConnectionSetupStatusResponse struct {
	IsCompleted bool                            `json:"isCompleted"`
	SetupStatus *enums.LLMConnectionSetupStatus `json:"setupStatus"`
}
