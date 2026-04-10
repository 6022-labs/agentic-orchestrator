package enums

import "strings"

type AgentContextMessageRole string
type AgentContextMessageRoleList []AgentContextMessageRole

const (
	AgentContextMessageRoleUser      AgentContextMessageRole = "user"
	AgentContextMessageRoleSystem    AgentContextMessageRole = "system"
	AgentContextMessageRoleAssistant AgentContextMessageRole = "assistant"
)

func GetAllAgentContextMessageRoles() AgentContextMessageRoleList {
	return AgentContextMessageRoleList{
		AgentContextMessageRoleUser,
		AgentContextMessageRoleSystem,
		AgentContextMessageRoleAssistant,
	}
}

func (k AgentContextMessageRole) String() string {
	return string(k)
}

func (k AgentContextMessageRoleList) String() string {
	var types []string
	for _, v := range k {
		types = append(types, string(v))
	}
	return strings.Join(types, ", ")
}

func (k AgentContextMessageRoleList) Contains(messageRole AgentContextMessageRole) bool {
	for _, v := range k {
		if v == messageRole {
			return true
		}
	}
	return false
}
