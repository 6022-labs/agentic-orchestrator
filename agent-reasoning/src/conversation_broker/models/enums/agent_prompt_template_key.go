package enums

import (
	"fmt"
	"strings"
)

type AgentPromptTemplateKey string
type AgentPromptTemplateKeyList []AgentPromptTemplateKey

const (
	AgentPromptTemplateKeyCloneSystem  AgentPromptTemplateKey = "clone-system"
	AgentPromptTemplateKeyCloneDefault AgentPromptTemplateKey = "clone-default"

	AgentPromptTemplateKeyFacilitatorSystem   AgentPromptTemplateKey = "facilitator-system"
	AgentPromptTemplateKeyFacilitatorDefault  AgentPromptTemplateKey = "facilitator-default"
	AgentPromptTemplateKeyFacilitatorConclude AgentPromptTemplateKey = "facilitator-conclude"

	AgentPromptTemplateKeyExpertSystem  AgentPromptTemplateKey = "expert-system"
	AgentPromptTemplateKeyExpertDefault AgentPromptTemplateKey = "expert-default"
)

func GetAllAgentPromptTemplateKeys() AgentPromptTemplateKeyList {
	return AgentPromptTemplateKeyList{
		AgentPromptTemplateKeyCloneSystem,
		AgentPromptTemplateKeyCloneDefault,
		AgentPromptTemplateKeyFacilitatorSystem,
		AgentPromptTemplateKeyFacilitatorDefault,
		AgentPromptTemplateKeyFacilitatorConclude,
		AgentPromptTemplateKeyExpertSystem,
		AgentPromptTemplateKeyExpertDefault,
	}
}

func (k AgentPromptTemplateKey) String() string {
	return string(k)
}

func (k AgentPromptTemplateKeyList) String() string {
	var keys []string
	for _, v := range k {
		keys = append(keys, string(v))
	}
	return strings.Join(keys, ", ")
}

func (k *AgentPromptTemplateKey) Set(value string) error {
	normalizedValue := strings.ToLower(value)
	for _, v := range GetAllAgentPromptTemplateKeys() {
		if normalizedValue == v.String() {
			*k = v
			return nil
		}
	}
	return fmt.Errorf("invalid AgentPromptTemplateKey: %s", value)
}

func (k AgentPromptTemplateKeyList) Contains(key AgentPromptTemplateKey) bool {
	for _, v := range k {
		if v == key {
			return true
		}
	}
	return false
}
