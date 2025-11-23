package prompt_parameters

import (
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type PromptParametersBuilder struct {
	agentRepository repositories.AgentRepositoryInterface
	params          map[string]any
}

func NewPromptParametersBuilder(agentRepository repositories.AgentRepositoryInterface) *PromptParametersBuilder {
	return &PromptParametersBuilder{
		agentRepository: agentRepository,
		params:          map[string]any{},
	}
}

func (b *PromptParametersBuilder) WithBase(base map[string]any) interfaces.PromptParametersBuilderInterface {
	out := make(map[string]any, len(base))
	for k, v := range base {
		out[k] = v
	}
	b.params = out
	return b
}

func (b *PromptParametersBuilder) WithParam(path []string, value any, parameterType *interfaces.PromptParameterType) interfaces.PromptParametersBuilderInterface {
	if len(path) == 0 {
		return b
	}

	m := b.params
	// Walk/create all but the last segment
	for _, key := range path[:len(path)-1] {
		if next, ok := m[key].(map[string]any); ok && next != nil {
			m = next
			continue
		}
		newMap := make(map[string]any)
		m[key] = newMap
		m = newMap
	}

	// Set final value
	lastKey := path[len(path)-1]
	m[lastKey] = value
	if parameterType != nil && *parameterType != "" {
		m[lastKey+".type"] = string(*parameterType)
	}

	return b
}

func (b *PromptParametersBuilder) Build() (map[string]any, error) {
	if err := b.applyTypeDirectives(); err != nil {
		return nil, err
	}

	out := make(map[string]any, len(b.params))
	for k, v := range b.params {
		out[k] = v
	}
	return out, nil
}

type typeDirective struct {
	Type interfaces.PromptParameterType
}

func (b *PromptParametersBuilder) applyTypeDirectives() error {
	return b.applyTypeDirectivesInValue(b.params)
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
			case interfaces.PromptParameterTypeAgent:
				if err := b.resolveAgentParameter(m, baseKey, val); err != nil {
					return err
				}
			case interfaces.PromptParameterTypeAgentArray:
				if err := b.resolveAgentArrayParameter(m, baseKey, val); err != nil {
					return err
				}
			case interfaces.PromptParameterTypePotentialAgent:
				if err := b.resolvePotentialAgentParameter(m, baseKey, val); err != nil {
					return err
				}
			case interfaces.PromptParameterTypePotentialAgentArray:
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

	agent, err := b.agentRepository.GetByEnsDomain(strVal)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get agent by ENS domain: %w", err)
	}
	if agent == nil {
		return fmt.Errorf("agent not found")
	}

	target[baseKey] = newAgentPromptParameterFromModel(agent)

	return nil
}

func (b *PromptParametersBuilder) resolveAgentArrayParameter(target map[string]any, baseKey string, value any) error {
	strVals, err := interfaceToStringSlice(value)
	if err != nil {
		return err
	}

	strVals = lo.Uniq(strVals)

	agents, err := b.agentRepository.GetByEnsDomains(strVals)
	if err != nil {
		return fmt.Errorf("failed to get agent by ENS domains: %w", err)
	}

	if len(agents) != len(strVals) {
		return fmt.Errorf("not all agents found")
	}

	agentParameters := make([]*agentPromptParameter, len(agents))

	for i, agent := range agents {
		agentParameters[i] = newAgentPromptParameterFromModel(agent)
	}

	target[baseKey] = agentParameters

	return nil
}

func (b *PromptParametersBuilder) resolvePotentialAgentParameter(target map[string]any, baseKey string, value any) error {
	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string value for potential agent, got %T", value)
	}

	agent, err := b.agentRepository.GetByEnsDomain(strVal)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get potential agent by ENS domain: %w", err)
	}

	target[baseKey] = newPotentialAgentPromptParameterFromModel(strVal, agent)

	return nil
}

func (b *PromptParametersBuilder) resolvePotentialAgentArrayParameter(target map[string]any, baseKey string, value any) error {
	strVals, err := interfaceToStringSlice(value)
	if err != nil {
		return err
	}

	strVals = lo.Uniq(strVals)

	agents, err := b.agentRepository.GetByEnsDomains(strVals)
	if err != nil {
		return fmt.Errorf("failed to get potential agents by ENS domains: %w", err)
	}

	potentialAgentParameters := make([]*potentialAgentPromptParameter, len(strVals))

	for i, strVal := range strVals {
		agent, _ := lo.Find(agents, func(a *models.Agent) bool {
			return a.EnsDomain == strVal
		})
		potentialAgentParameters[i] = newPotentialAgentPromptParameterFromModel(strVal, agent)
	}

	target[baseKey] = potentialAgentParameters

	return nil
}

func interfaceToStringSlice(vals any) ([]string, error) {
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
