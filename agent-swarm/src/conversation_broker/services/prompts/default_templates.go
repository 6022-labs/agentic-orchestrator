package prompts

import "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"

// defaultTemplates are the built-in fallbacks used when an agent has no
// stored override for a key; per-agent personality lives in the agent node.
var defaultTemplates = map[enums.AgentPromptTemplateKey]string{
	enums.AgentPromptTemplateKeyCloneSystem: `You are taking part in a group conversation` +
		`{{if .ChannelDescription}} about: {{.ChannelDescription}}{{end}}.
Participants: {{range .ChannelMembers}}{{.Name}} {{end}}.
Answer as yourself, stay concise and only address what concerns you.`,

	enums.AgentPromptTemplateKeyCloneDefault: `{{.LatestRequest.Author.Name}} said: {{.LatestRequest.Content}}
Reply to the conversation.`,

	enums.AgentPromptTemplateKeyExpertSystem: `You are a domain expert taking part in a group conversation` +
		`{{if .ChannelDescription}} about: {{.ChannelDescription}}{{end}}.
Participants: {{range .ChannelMembers}}{{.Name}} {{end}}.
Provide precise, factual answers within your expertise.`,

	enums.AgentPromptTemplateKeyExpertDefault: `{{.LatestRequest.Author.Name}} said: {{.LatestRequest.Content}}
Answer with your expert view.`,

	enums.AgentPromptTemplateKeyFacilitatorSystem: `You facilitate a group conversation` +
		`{{if .ChannelDescription}} about: {{.ChannelDescription}}{{end}}.
Participants: {{range .ChannelMembers}}{{.Name}} {{end}}.
Keep the discussion moving toward the goal stated in the initial request: {{.InitialRequest.Content}}
When you judge the conversation complete, include the marker :end: in your reply.`,

	enums.AgentPromptTemplateKeyFacilitatorDefault: `The conversation now counts ` +
		`{{.ConversationMessageCount}} messages. Latest from {{.LatestRequest.Author.Name}}: {{.LatestRequest.Content}}
Facilitate the next step.`,

	enums.AgentPromptTemplateKeyFacilitatorConclude: `The conversation is ending now. ` +
		`Write a short conclusion summarizing the outcome for all participants.`,
}

func (p *promptTemplateProvider) defaultTemplateForKey(key enums.AgentPromptTemplateKey) (string, bool) {
	tpl, ok := defaultTemplates[key]
	return tpl, ok
}
