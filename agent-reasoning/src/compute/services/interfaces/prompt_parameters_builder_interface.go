package interfaces

type PromptParameterType string
type PromptParameterTypeList []PromptParameterType

const (
	PromptParameterTypeAgent               PromptParameterType = "agent"
	PromptParameterTypeAgentArray          PromptParameterType = "agent[]"
	PromptParameterTypePotentialAgent      PromptParameterType = "potential_agent"
	PromptParameterTypePotentialAgentArray PromptParameterType = "potential_agent[]"
)

func GetAllPromptParameterTypes() PromptParameterTypeList {
	return PromptParameterTypeList{
		PromptParameterTypeAgent,
		PromptParameterTypeAgentArray,
		PromptParameterTypePotentialAgent,
		PromptParameterTypePotentialAgentArray,
	}
}

func (t PromptParameterTypeList) Contains(target PromptParameterType) bool {
	for _, item := range t {
		if item == target {
			return true
		}
	}
	return false
}

func (t PromptParameterType) IsValid() bool {
	return GetAllPromptParameterTypes().Contains(t)
}

func (t PromptParameterType) String() string {
	return string(t)
}

func (t PromptParameterType) ToPointer() *PromptParameterType {
	return &t
}

type PromptParametersBuilderInterface interface {
	WithBase(base map[string]any) PromptParametersBuilderInterface
	WithParam(path []string, value any, parameterType *PromptParameterType) PromptParametersBuilderInterface
	Build() (map[string]any, error)
}
