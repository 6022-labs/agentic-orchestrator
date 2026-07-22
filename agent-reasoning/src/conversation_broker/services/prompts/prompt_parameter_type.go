package prompts

import (
	"fmt"
	"strings"
)

type PromptParameterType string

const (
	PromptParameterTypeAgent               PromptParameterType = "agent"
	PromptParameterTypeAgentArray          PromptParameterType = "agent[]"
	PromptParameterTypePotentialAgent      PromptParameterType = "potential_agent"
	PromptParameterTypePotentialAgentArray PromptParameterType = "potential_agent[]"
)

func (t PromptParameterType) String() string {
	return string(t)
}

func PromptParameterTypeFromString(s string) (PromptParameterType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "agent":
		return PromptParameterTypeAgent, nil
	case "agents", "agent[]":
		return PromptParameterTypeAgentArray, nil
	case "potential_agent":
		return PromptParameterTypePotentialAgent, nil
	case "potential_agents", "potential_agent[]":
		return PromptParameterTypePotentialAgentArray, nil
	default:
		return "", fmt.Errorf("unknown prompt parameter type: %s", s)
	}
}
