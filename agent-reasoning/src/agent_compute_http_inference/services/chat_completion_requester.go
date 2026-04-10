package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	orchestrator_agent_compute_http_inference_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/settings"
	"github.com/pkoukk/tiktoken-go"
	"go.uber.org/zap"
)

type ChatCompletionRequester struct {
	logger *zap.Logger

	tokenSettings *settings.TokenSettings

	xSettings             *settings.XSettings
	openAiSettings        *settings.OpenAISettings
	ollamaSettings        *settings.OllamaSettings
	mistralSettings       *settings.MistralSettings
	deepSeekSettings      *settings.DeepSeekSettings
	libertAiSettings      *settings.LibertAISettings
	anthropicSettings     *settings.AnthropicSettings
	amazonBedrockSettings *settings.AmazonBedrockSettings

	openAIResponsesStandardChatCompletionRequester       OpenAIResponsesStandardChatCompletionRequesterInterface
	ollamaChatMessageStandardChatCompletionRequester     OllamaChatMessageStandardChatCompletionRequesterInterface
	anthropicMessagesStandardChatCompletionRequester     AnthropicMessagesStandardChatCompletionRequesterInterface
	openAICompletionsStandardChatCompletionRequester     OpenAICompletionsStandardChatCompletionRequesterInterface
	amazonBedrockConverseStandardChatCompletionRequester AmazonBedrockConverseStandardChatCompletionRequesterInterface
}

func NewChatCompletionRequester(
	logger *zap.Logger,
	xSettings *settings.XSettings,
	tokenSettings *settings.TokenSettings,
	openAiSettings *settings.OpenAISettings,
	ollamaSettings *settings.OllamaSettings,
	mistralSettings *settings.MistralSettings,
	deepSeekSettings *settings.DeepSeekSettings,
	libertAiSettings *settings.LibertAISettings,
	anthropicSettings *settings.AnthropicSettings,
	amazonBedrockSettings *settings.AmazonBedrockSettings,
	openAIResponsesStandardChatCompletionRequester OpenAIResponsesStandardChatCompletionRequesterInterface,
	ollamaChatMessageStandardChatCompletionRequester OllamaChatMessageStandardChatCompletionRequesterInterface,
	anthropicMessagesStandardChatCompletionRequester AnthropicMessagesStandardChatCompletionRequesterInterface,
	openAICompletionsStandardChatCompletionRequester OpenAICompletionsStandardChatCompletionRequesterInterface,
	amazonBedrockConverseStandardChatCompletionRequester AmazonBedrockConverseStandardChatCompletionRequesterInterface,
) *ChatCompletionRequester {
	return &ChatCompletionRequester{
		logger: logger,

		tokenSettings: tokenSettings,

		xSettings:             xSettings,
		openAiSettings:        openAiSettings,
		ollamaSettings:        ollamaSettings,
		mistralSettings:       mistralSettings,
		deepSeekSettings:      deepSeekSettings,
		libertAiSettings:      libertAiSettings,
		anthropicSettings:     anthropicSettings,
		amazonBedrockSettings: amazonBedrockSettings,

		openAIResponsesStandardChatCompletionRequester:       openAIResponsesStandardChatCompletionRequester,
		ollamaChatMessageStandardChatCompletionRequester:     ollamaChatMessageStandardChatCompletionRequester,
		openAICompletionsStandardChatCompletionRequester:     openAICompletionsStandardChatCompletionRequester,
		anthropicMessagesStandardChatCompletionRequester:     anthropicMessagesStandardChatCompletionRequester,
		amazonBedrockConverseStandardChatCompletionRequester: amazonBedrockConverseStandardChatCompletionRequester,
	}
}

func (r *ChatCompletionRequester) ChatCompletion(
	llmConnection *models.LLMConnection,
	mcpConnections []*models.MCPConnection,
	messages []*models.ChatCompletionMessage,
	apiKey *string,
) (*string, error) {
	trimmedMessages, err := r.trimMessagesToTokenLimit(messages, r.tokenSettings.MaxTokensPerRequest)
	if err != nil {
		return nil, err
	}

	platformSettings, err := r.getPlatformSettings(llmConnection.Platform)
	if err != nil {
		return nil, err
	}

	return r.completeUsingAppropriateStandard(
		llmConnection,
		mcpConnections,
		trimmedMessages,
		platformSettings,
		apiKey,
	)
}

// trimMessagesToTokenLimit trims messages from the start until the total token count is <= limit
func (r *ChatCompletionRequester) trimMessagesToTokenLimit(messages []*models.ChatCompletionMessage, limit uint) ([]*models.ChatCompletionMessage, error) {
	tokenizer, err := tiktoken.EncodingForModel("gpt-4")
	if err != nil {
		return nil, fmt.Errorf("failed to get tokenizer: %w", err)
	}

	totalTokens := uint(0)
	msgTokens := make([]uint, len(messages))
	for i, msg := range messages {
		tokens := tokenizer.Encode(msg.Content, nil, nil)
		msgTokens[i] = uint(len(tokens))
		totalTokens += uint(len(tokens))
	}

	startIdx := 0
	for totalTokens > limit && startIdx < len(messages) {
		totalTokens -= msgTokens[startIdx]
		startIdx++
	}

	return messages[startIdx:], nil
}

func (r *ChatCompletionRequester) getPlatformSettings(platform enums.LLMConnectionPlatform) (settings.StandardSettingsInterface, error) {
	switch platform {
	case enums.LLMConnectionPlatformLibertAI:
		return r.libertAiSettings, nil
	case enums.LLMConnectionPlatformOpenAI:
		return r.openAiSettings, nil
	case enums.LLMConnectionPlatformOllama:
		return r.ollamaSettings, nil
	case enums.LLMConnectionPlatformMistral:
		return r.mistralSettings, nil
	case enums.LLMConnectionPlatformDeepseek:
		return r.deepSeekSettings, nil
	case enums.LLMConnectionPlatformX:
		return r.xSettings, nil
	case enums.LLMConnectionPlatformAnthropic:
		return r.anthropicSettings, nil
	case enums.LLMConnectionPlatformAmazonBedrock:
		return r.amazonBedrockSettings, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}

func (r *ChatCompletionRequester) completeUsingAppropriateStandard(
	llmConnection *models.LLMConnection,
	mcpConnections []*models.MCPConnection,
	messages []*models.ChatCompletionMessage,
	platformSettings settings.StandardSettingsInterface,
	apiKey *string,
) (*string, error) {
	if platformSettings.PlatformAPI() == settings.PlatformAPIOpenAIResponsesStandard {
		inputs := orchestrator_agent_compute_http_inference_models.NewExternalOpenAICompletionsStandardRequestInputsFromChatCompletionMessages(messages)

		if apiKey == nil {
			return nil, fmt.Errorf("api key is required for OpenAI Responses Standard")
		}

		return r.openAIResponsesStandardChatCompletionRequester.ChatCompletion(
			platformSettings.GetBaseUrl(),
			*apiKey,
			platformSettings.HasDirectAccessToMcp(),
			llmConnection,
			inputs,
			mcpConnections,
			nil,
		)
	}

	if platformSettings.PlatformAPI() == settings.PlatformAPIOpenAICompletionsStandard {
		inputs := orchestrator_agent_compute_http_inference_models.NewExternalOpenAICompletionsStandardRequestMessagesFromChatCompletionMessages(messages)

		if apiKey == nil {
			return nil, fmt.Errorf("api key is required for OpenAI Completions Standard")
		}

		return r.openAICompletionsStandardChatCompletionRequester.ChatCompletion(
			platformSettings.GetBaseUrl(),
			*apiKey,
			llmConnection,
			inputs,
			mcpConnections,
		)
	}

	if platformSettings.PlatformAPI() == settings.PlatformAPIOllamaChatMessageStandard {
		inputs := orchestrator_agent_compute_http_inference_models.NewExternalOllamaChatMessageStandardMessagesFromChatCompletionMessages(messages)

		if llmConnection.BaseUrl == nil {
			return nil, fmt.Errorf("base URL is required for Ollama Chat Message Standard")
		}

		return r.ollamaChatMessageStandardChatCompletionRequester.ChatCompletion(
			*llmConnection.BaseUrl,
			llmConnection,
			inputs,
			mcpConnections,
		)
	}

	if platformSettings.PlatformAPI() == settings.PlatformAPIAnthropicMessagesStandard {
		inputs := orchestrator_agent_compute_http_inference_models.NewExternalAnthropicMessagesStandardMessagesFromChatCompletionMessages(messages)
		systems := orchestrator_agent_compute_http_inference_models.NewExternalAnthropicMessagesStandardRequestSystemsFromChatCompletionMessages(messages)

		if apiKey == nil {
			return nil, fmt.Errorf("api key is required for Anthropic Messages Standard")
		}

		return r.anthropicMessagesStandardChatCompletionRequester.ChatCompletion(
			platformSettings.GetBaseUrl(),
			*apiKey,
			platformSettings.HasDirectAccessToMcp(),
			llmConnection,
			inputs,
			systems,
			mcpConnections,
		)
	}

	if platformSettings.PlatformAPI() == settings.PlatformAPIAmazonBedrockConverseStandard {
		messages, messagesSystems := orchestrator_agent_compute_http_inference_models.NewExternalAmazonBedrockConverseStandardRequestMessagesAndSystemFromChatCompletionMessages(messages)

		if apiKey == nil {
			return nil, fmt.Errorf("api key is required for Amazon Bedrock Converse Standard")
		}
		if llmConnection.BaseUrl == nil {
			return nil, fmt.Errorf("base URL is required for Amazon Bedrock Converse Standard")
		}

		return r.amazonBedrockConverseStandardChatCompletionRequester.ChatCompletion(
			*llmConnection.BaseUrl,
			*apiKey,
			llmConnection,
			messages,
			messagesSystems,
			mcpConnections,
		)
	}

	return nil, fmt.Errorf("unsupported platform API: %s", platformSettings.PlatformAPI())
}
