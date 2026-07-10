package enums

type RemoteAgentInputRole string

const (
	RemoteAgentInputRoleUser      RemoteAgentInputRole = "user"
	RemoteAgentInputRoleAssistant RemoteAgentInputRole = "assistant"
)

func (r RemoteAgentInputRole) String() string {
	return string(r)
}
