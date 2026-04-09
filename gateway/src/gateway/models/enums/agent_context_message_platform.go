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

func GetAllAgentContextMessagePlatformWithoutNone() AgentContextMessagePlatformList {
	return AgentContextMessagePlatformList{
		AgentContextMessagePlatformSlack,
	}
}

func (k AgentContextMessagePlatformList) String() string {
	var types []string
	for _, v := range k {
		types = append(types, string(v))
	}
	return strings.Join(types, ", ")
}

func (list AgentContextMessagePlatformList) Contains(value AgentContextMessagePlatform) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
