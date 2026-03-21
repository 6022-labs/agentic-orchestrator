package prompt_parameters

import (
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
)

func PromptParameterTypeFromString(s string) (interfaces.PromptParameterType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "agent":
		return interfaces.PromptParameterTypeAgent, nil
	case "agents", "agent[]":
		return interfaces.PromptParameterTypeAgentArray, nil
	case "potential_agent":
		return interfaces.PromptParameterTypePotentialAgent, nil
	case "potential_agents", "potential_agent[]":
		return interfaces.PromptParameterTypePotentialAgentArray, nil
	default:
		return "", fmt.Errorf("unknown prompt parameter type: %s", s)
	}
}
