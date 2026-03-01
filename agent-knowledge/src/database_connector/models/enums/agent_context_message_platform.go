package enums

import "strings"

type AgentContextMessagePlatform string
type AgentContextMessagePlatformList []AgentContextMessagePlatform

const (
	AgentContextMessagePlatformNone  AgentContextMessagePlatform = "none"
	AgentContextMessagePlatformSlack AgentContextMessagePlatform = "slack"
)

func GetAllAgentContextMessagePlatforms() AgentContextMessagePlatformList {
	return AgentContextMessagePlatformList{
		AgentContextMessagePlatformNone,
		AgentContextMessagePlatformSlack,
	}
}

func (k AgentContextMessagePlatform) String() string {
	return string(k)
}

func (k AgentContextMessagePlatformList) String() string {
	var types []string
	for _, v := range k {
		types = append(types, string(v))
	}
	return strings.Join(types, ", ")
}

func (k AgentContextMessagePlatformList) Contains(messagePlatform AgentContextMessagePlatform) bool {
	for _, v := range k {
		if v == messagePlatform {
			return true
		}
	}
	return false
}
