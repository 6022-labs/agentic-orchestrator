package models

import "encoding/json"

type ExternalOpenAICompletionsStandardToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

func (f *ExternalOpenAICompletionsStandardToolCallFunction) ParsedArguments() (map[string]any, error) {
	var args map[string]any
	err := json.Unmarshal([]byte(f.Arguments), &args)
	if err != nil {
		return nil, err
	}

	return args, nil
}
