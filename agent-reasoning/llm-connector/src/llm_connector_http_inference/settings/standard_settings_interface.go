package settings

type PlatformAPIInterface string

const (
	PlatformAPIOpenAIResponsesStandard       PlatformAPIInterface = "OpenAIResponsesStandard"
	PlatformAPIOpenAICompletionsStandard     PlatformAPIInterface = "OpenAICompletionsStandard"
	PlatformAPIOllamaChatMessageStandard     PlatformAPIInterface = "OllamaChatMessageStandard"
	PlatformAPIAnthropicMessagesStandard     PlatformAPIInterface = "AnthropicMessagesStandard"
	PlatformAPIAmazonBedrockConverseStandard PlatformAPIInterface = "AmazonBedrockConverseStandard"
)

type StandardSettingsInterface interface {
	GetBaseUrl() string
	HasDirectAccessToMcp() bool
	PlatformAPI() PlatformAPIInterface
}
