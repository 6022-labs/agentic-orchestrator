package prompts

import (
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/samber/lo"
)

// PromptParametersBuilder assembles the template parameter map, resolving
// ".type"-suffixed directives (agent / potential_agent lookups) through the
// registry port.
type PromptParametersBuilder struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
	params                 map[string]any
}

func NewPromptParametersBuilder(agentRegistryRequester interfaces.AgentRegistryRequesterInterface) *PromptParametersBuilder {
	return &PromptParametersBuilder{
		agentRegistryRequester: agentRegistryRequester,
		params:                 map[string]any{},
	}
}

func (b *PromptParametersBuilder) WithBase(base map[string]any) *PromptParametersBuilder {
	out := make(map[string]any, len(base))
	for k, v := range base {
		out[k] = v
	}
	b.params = out
	return b
}

func (b *PromptParametersBuilder) WithParam(path []string, value any, parameterType *PromptParameterType) *PromptParametersBuilder {
	if len(path) == 0 {
		return b
	}

	m := b.params
	for _, key := range path[:len(path)-1] {
		if next, ok := m[key].(map[string]any); ok && next != nil {
			m = next
			continue
		}
		newMap := make(map[string]any)
		m[key] = newMap
		m = newMap
	}

	lastKey := path[len(path)-1]
	m[lastKey] = value
	if parameterType != nil && *parameterType != "" {
		m[lastKey+".type"] = string(*parameterType)
	}

	return b
}

func (b *PromptParametersBuilder) Build() (map[string]any, error) {
	if err := b.applyTypeDirectivesInValue(b.params); err != nil {
		return nil, err
	}

	out := make(map[string]any, len(b.params))
	for k, v := range b.params {
		out[k] = v
	}
	return out, nil
}

type typeDirective struct {
	Type PromptParameterType
}

func (b *PromptParametersBuilder) applyTypeDirectivesInValue(v any) error {
	switch vv := v.(type) {
	case map[string]any:
		if err := b.applyTypeDirectivesInMap(vv); err != nil {
			return err
		}
		for _, child := range vv {
			if err := b.applyTypeDirectivesInValue(child); err != nil {
				return err
			}
		}
	case []any:
		for _, item := range vv {
			if err := b.applyTypeDirectivesInValue(item); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *PromptParametersBuilder) applyTypeDirectivesInMap(m map[string]any) error {
	dotted := map[string]*typeDirective{}
	for k, v := range m {
		if strings.HasSuffix(k, ".type") {
			baseKey := strings.TrimSuffix(k, ".type")
			t, _ := v.(string)
			promptParameterType, err := PromptParameterTypeFromString(t)
			if err != nil {
				return fmt.Errorf("failed to parse prompt parameter type: %w", err)
			}
			dotted[baseKey] = &typeDirective{Type: promptParameterType}
		}
	}

	for baseKey, spec := range dotted {
		if spec == nil || spec.Type == "" {
			continue
		}
		val, ok := m[baseKey]
		if ok {
			switch spec.Type {
			case PromptParameterTypeAgent:
				if err := b.resolveAgentParameter(m, baseKey, val); err != nil {
					return err
				}
			case PromptParameterTypeAgentArray:
				if err := b.resolveAgentArrayParameter(m, baseKey, val); err != nil {
					return err
				}
			case PromptParameterTypePotentialAgent:
				if err := b.resolvePotentialAgentParameter(m, baseKey, val); err != nil {
					return err
				}
			case PromptParameterTypePotentialAgentArray:
				if err := b.resolvePotentialAgentArrayParameter(m, baseKey, val); err != nil {
					return err
				}
			}
		}

		delete(m, baseKey+".type")
	}

	return nil
}

func (b *PromptParametersBuilder) resolveAgentParameter(target map[string]any, baseKey string, value any) error {
	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string value for agent, got %T", value)
	}

	agent, err := b.agentRegistryRequester.GetMintedAgentByEnsDomain(strVal)
	if err != nil {
		return fmt.Errorf("failed to get agent by ENS domain: %w", err)
	}
	if agent == nil {
		return fmt.Errorf("agent not found")
	}

	target[baseKey] = newAgentPromptParameter(agent)

	return nil
}

func (b *PromptParametersBuilder) resolveAgentArrayParameter(target map[string]any, baseKey string, value any) error {
	strVals, err := b.interfaceToStringSlice(value)
	if err != nil {
		return err
	}

	strVals = lo.Uniq(strVals)

	agentParameters := make([]*agentPromptParameter, len(strVals))
	for i, ensDomain := range strVals {
		agent, err := b.agentRegistryRequester.GetMintedAgentByEnsDomain(ensDomain)
		if err != nil {
			return fmt.Errorf("failed to get agents by ENS domains: %w", err)
		}
		if agent == nil {
			return fmt.Errorf("not all agents found")
		}
		agentParameters[i] = newAgentPromptParameter(agent)
	}

	target[baseKey] = agentParameters

	return nil
}

func (b *PromptParametersBuilder) resolvePotentialAgentParameter(target map[string]any, baseKey string, value any) error {
	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string value for potential agent, got %T", value)
	}

	agent, err := b.agentRegistryRequester.GetMintedAgentByEnsDomain(strVal)
	if err != nil {
		return fmt.Errorf("failed to get potential agent by ENS domain: %w", err)
	}

	target[baseKey] = newPotentialAgentPromptParameter(strVal, agent)

	return nil
}

func (b *PromptParametersBuilder) resolvePotentialAgentArrayParameter(target map[string]any, baseKey string, value any) error {
	strVals, err := b.interfaceToStringSlice(value)
	if err != nil {
		return err
	}

	strVals = lo.Uniq(strVals)

	agentMap := make(map[string]*models.RegistryAgent, len(strVals))
	for _, ensDomain := range strVals {
		agent, err := b.agentRegistryRequester.GetMintedAgentByEnsDomain(ensDomain)
		if err != nil {
			return fmt.Errorf("failed to get potential agents by ENS domains: %w", err)
		}
		if agent != nil {
			agentMap[ensDomain] = agent
		}
	}

	potentialAgentParameters := make([]*potentialAgentPromptParameter, len(strVals))

	for i, strVal := range strVals {
		potentialAgentParameters[i] = newPotentialAgentPromptParameter(strVal, agentMap[strVal])
	}

	target[baseKey] = potentialAgentParameters

	return nil
}

func (b *PromptParametersBuilder) interfaceToStringSlice(vals any) ([]string, error) {
	out := []string{}

	switch v := vals.(type) {
	case []string:
		out = v
	case []any:
		for index, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("expected potential agent at index %d to be string, got %T", index, item)
			}
			out = append(out, s)
		}
	default:
		return nil, fmt.Errorf("unsupported type: %T", vals)
	}

	return out, nil
}
