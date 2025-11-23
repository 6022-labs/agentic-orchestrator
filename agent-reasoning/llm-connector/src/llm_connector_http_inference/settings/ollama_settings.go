package settings

type OllamaSettings struct {
}

func NewOllamaSettings() *OllamaSettings {
	return &OllamaSettings{}
}

func (s *OllamaSettings) GetBaseUrl() string {
	// Not used for Ollama, only here to satisfy the interface
	return ""
}

func (s *OllamaSettings) HasDirectAccessToMcp() bool {
	return false
}

func (s *OllamaSettings) PlatformAPI() PlatformAPIInterface {
	return PlatformAPIOllamaChatMessageStandard
}
