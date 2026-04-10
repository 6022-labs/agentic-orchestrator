package enums

import "strings"

type AgentPromptTemplateKey string
type AgentPromptTemplateKeyList []AgentPromptTemplateKey

const (
	AgentPromptTemplateKeyCloneSystem  AgentPromptTemplateKey = "clone-system"
	AgentPromptTemplateKeyCloneDefault AgentPromptTemplateKey = "clone-default"

	AgentPromptTemplateKeyFacilitatorSystem   AgentPromptTemplateKey = "facilitator-system"
	AgentPromptTemplateKeyFacilitatorDefault  AgentPromptTemplateKey = "facilitator-default"
	AgentPromptTemplateKeyFacilitatorConclude AgentPromptTemplateKey = "facilitator-conclude"

	AgentPromptTemplateKeyExpertSystem          AgentPromptTemplateKey = "expert-system"
	AgentPromptTemplateKeyExpertContextFound    AgentPromptTemplateKey = "expert-context-found"
	AgentPromptTemplateKeyExpertContextNotFound AgentPromptTemplateKey = "expert-context-not-found"
)

func GetAllAgentPromptTemplateKeys() AgentPromptTemplateKeyList {
	return AgentPromptTemplateKeyList{
		AgentPromptTemplateKeyCloneSystem,
		AgentPromptTemplateKeyCloneDefault,
		AgentPromptTemplateKeyFacilitatorSystem,
		AgentPromptTemplateKeyFacilitatorDefault,
		AgentPromptTemplateKeyFacilitatorConclude,
		AgentPromptTemplateKeyExpertSystem,
		AgentPromptTemplateKeyExpertContextFound,
		AgentPromptTemplateKeyExpertContextNotFound,
	}
}

func GetAllExpertAgentPromptTemplateKeys() AgentPromptTemplateKeyList {
	return AgentPromptTemplateKeyList{
		AgentPromptTemplateKeyExpertSystem,
		AgentPromptTemplateKeyExpertContextFound,
		AgentPromptTemplateKeyExpertContextNotFound,
	}
}

func GetAllFacilitatorAgentPromptTemplateKeys() AgentPromptTemplateKeyList {
	return AgentPromptTemplateKeyList{
		AgentPromptTemplateKeyFacilitatorSystem,
		AgentPromptTemplateKeyFacilitatorDefault,
		AgentPromptTemplateKeyFacilitatorConclude,
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

func (k AgentPromptTemplateKeyList) Contains(key AgentPromptTemplateKey) bool {
	for _, v := range k {
		if v == key {
			return true
		}
	}
	return false
}
