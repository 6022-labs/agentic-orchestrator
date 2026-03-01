package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models/enums"
	external_models "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/llm_connector_http_inference_mocks/services_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenRequestingChatCompletingTestingSuite struct {
	sut *services.ChatCompletionRequester

	tokenSettings *settings.TokenSettings

	xSettings             *settings.XSettings
	openAiSettings        *settings.OpenAISettings
	mistralSettings       *settings.MistralSettings
	deepSeekSettings      *settings.DeepSeekSettings
	libertAiSettings      *settings.LibertAISettings
	anthropicSettings     *settings.AnthropicSettings
	amazonBedrockSettings *settings.AmazonBedrockSettings

	mockOpenAIResponsesStandardChatCompletionRequester       services_mocks.MockOpenAIResponsesStandardChatCompletionRequesterInterface
	mockOllamaChatMessageStandardChatCompletionRequester     services_mocks.MockOllamaChatMessageStandardChatCompletionRequesterInterface
	mockAnthropicMessagesStandardChatCompletionRequester     services_mocks.MockAnthropicMessagesStandardChatCompletionRequesterInterface
	mockOpenAICompletionsStandardChatCompletionRequester     services_mocks.MockOpenAICompletionsStandardChatCompletionRequesterInterface
	mockAmazonBedrockConverseStandardChatCompletionRequester services_mocks.MockAmazonBedrockConverseStandardChatCompletionRequesterInterface
}

func WhenRequestingChatCompletionBeforeEach(t *testing.T) *WhenRequestingChatCompletingTestingSuite {
	mockController := gomock.NewController(t)

	mockOpenAIResponsesStandardChatCompletionRequester := services_mocks.NewMockOpenAIResponsesStandardChatCompletionRequesterInterface(mockController)
	mockOllamaChatMessageStandardChatCompletionRequester := services_mocks.NewMockOllamaChatMessageStandardChatCompletionRequesterInterface(mockController)
	mockAnthropicMessagesStandardChatCompletionRequester := services_mocks.NewMockAnthropicMessagesStandardChatCompletionRequesterInterface(mockController)
	mockOpenAICompletionsStandardChatCompletionRequester := services_mocks.NewMockOpenAICompletionsStandardChatCompletionRequesterInterface(mockController)
	mockAmazonBedrockConverseStandardChatCompletionRequester := services_mocks.NewMockAmazonBedrockConverseStandardChatCompletionRequesterInterface(mockController)

	tokenSettings := &settings.TokenSettings{
		MaxTokensPerRequest: 4096,
	}

	xSettings := &settings.XSettings{}
	openAiSettings := &settings.OpenAISettings{}
	ollamaSettings := &settings.OllamaSettings{}
	mistralSettings := &settings.MistralSettings{}
	deepSeekSettings := &settings.DeepSeekSettings{}
	libertAiSettings := &settings.LibertAISettings{}
	anthropicSettings := &settings.AnthropicSettings{}
	amazonBedrockSettings := &settings.AmazonBedrockSettings{}

	sut := services.NewChatCompletionRequester(
		zap.NewNop(),
		xSettings,
		tokenSettings,
		openAiSettings,
		ollamaSettings,
		mistralSettings,
		deepSeekSettings,
		libertAiSettings,
		anthropicSettings,
		amazonBedrockSettings,
		mockOpenAIResponsesStandardChatCompletionRequester,
		mockOllamaChatMessageStandardChatCompletionRequester,
		mockAnthropicMessagesStandardChatCompletionRequester,
		mockOpenAICompletionsStandardChatCompletionRequester,
		mockAmazonBedrockConverseStandardChatCompletionRequester,
	)

	return &WhenRequestingChatCompletingTestingSuite{
		sut: sut,

		tokenSettings: tokenSettings,

		xSettings:             xSettings,
		openAiSettings:        openAiSettings,
		mistralSettings:       mistralSettings,
		deepSeekSettings:      deepSeekSettings,
		libertAiSettings:      libertAiSettings,
		anthropicSettings:     anthropicSettings,
		amazonBedrockSettings: amazonBedrockSettings,

		mockOpenAIResponsesStandardChatCompletionRequester:       *mockOpenAIResponsesStandardChatCompletionRequester,
		mockOllamaChatMessageStandardChatCompletionRequester:     *mockOllamaChatMessageStandardChatCompletionRequester,
		mockAnthropicMessagesStandardChatCompletionRequester:     *mockAnthropicMessagesStandardChatCompletionRequester,
		mockOpenAICompletionsStandardChatCompletionRequester:     *mockOpenAICompletionsStandardChatCompletionRequester,
		mockAmazonBedrockConverseStandardChatCompletionRequester: *mockAmazonBedrockConverseStandardChatCompletionRequester,
	}
}

func TestWhenRequestingChatCompletion(t *testing.T) {
	t.Parallel()

	t.Run("Given llm connection is LibertAI", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingChatCompletingTestingSuite, llmConnection *models.LLMConnection, messages []*models.ChatCompletionMessage, apiKey string) {
			suite.libertAiSettings.BaseUrl = "http://libertai.api"

			result := "completion-result"

			suite.mockOpenAICompletionsStandardChatCompletionRequester.
				EXPECT().
				ChatCompletion(
					suite.libertAiSettings.GetBaseUrl(),
					apiKey,
					llmConnection,
					gomock.Any(),
					nil,
				).
				DoAndReturn(func(_, _ string, _ *models.LLMConnection, externalMessages []external_models.ExternalOpenAICompletionsStandardRequestMessage, _ []*models.MCPConnection) (*string, error) {
					assert.Equal(t, len(messages), len(externalMessages))
					for i, message := range messages {
						assert.Equal(t, message.Role, externalMessages[i].Role)
						assert.Equal(t, message.Content, *externalMessages[i].Content)
					}

					return &result, nil
				})
		}

		t.Run("Should use OpenAI completion standard", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingChatCompletionBeforeEach(t)

			llmConnection := model_builders.NewLLMConnectionBuilder().
				WithPlatform(enums.LLMConnectionPlatformLibertAI).
				WithModel("gpt-4").
				Build()
			message := model_builders.NewChatCompletionMessageBuilder().Build()
			apiKey := "api-key"

			initSuite(suite, llmConnection, []*models.ChatCompletionMessage{message}, apiKey)

			result, err := suite.sut.ChatCompletion(llmConnection, nil, []*models.ChatCompletionMessage{message}, &apiKey)

			assert.Nil(t, err)
			assert.NotNil(t, result)
		})
	})

	t.Run("Given llm connection is OpenAI", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingChatCompletingTestingSuite, llmConnection *models.LLMConnection, messages []*models.ChatCompletionMessage, apiKey string) {
			suite.openAiSettings.BaseUrl = "http://openai.api"

			result := "completion-result"

			suite.mockOpenAIResponsesStandardChatCompletionRequester.
				EXPECT().
				ChatCompletion(
					suite.openAiSettings.GetBaseUrl(),
					apiKey,
					suite.openAiSettings.HasDirectAccessToMcp(),
					llmConnection,
					gomock.Any(),
					nil,
					nil,
				).
				DoAndReturn(func(_, _ string, _ bool, llmConnection *models.LLMConnection, externalMessages []external_models.ExternalOpenAIResponsesStandardRequestInput, _ []*models.MCPConnection, _ *string) (*string, error) {
					assert.Equal(t, len(messages), len(externalMessages))
					for i, message := range messages {
						assert.Equal(t, message.Role, *externalMessages[i].Role)
						assert.Equal(t, message.Content, *externalMessages[i].Content)
					}

					return &result, nil
				})
		}

		t.Run("Should use OpenAI responses standard", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingChatCompletionBeforeEach(t)

			llmConnection := model_builders.NewLLMConnectionBuilder().
				WithPlatform(enums.LLMConnectionPlatformOpenAI).
				WithModel("gpt-4").
				Build()
			message := model_builders.NewChatCompletionMessageBuilder().Build()
			apiKey := "api-key"

			initSuite(suite, llmConnection, []*models.ChatCompletionMessage{message}, apiKey)

			result, err := suite.sut.ChatCompletion(llmConnection, nil, []*models.ChatCompletionMessage{message}, &apiKey)

			assert.Nil(t, err)
			assert.NotNil(t, result)
		})
	})

	t.Run("Given llm connection is Ollama", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingChatCompletingTestingSuite, llmConnection *models.LLMConnection, messages []*models.ChatCompletionMessage) {
			result := "completion-result"

			suite.mockOllamaChatMessageStandardChatCompletionRequester.
				EXPECT().
				ChatCompletion(
					*llmConnection.BaseUrl,
					llmConnection,
					gomock.Any(),
					nil,
				).
				DoAndReturn(func(_ string, _ *models.LLMConnection, externalMessages []external_models.ExternalOllamaChatMessageStandardMessage, _ []*models.MCPConnection) (*string, error) {
					assert.Equal(t, len(messages), len(externalMessages))
					for i, message := range messages {
						assert.Equal(t, message.Role, externalMessages[i].Role)
						assert.Equal(t, message.Content, externalMessages[i].Content)
					}

					return &result, nil
				})
		}

		t.Run("Should use Ollama chat message standard", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingChatCompletionBeforeEach(t)

			llmConnection := model_builders.NewLLMConnectionBuilder().
				WithPlatform(enums.LLMConnectionPlatformOllama).
				WithBaseUrl("http://ollama.api").
				WithModel("ollama-model").
				Build()
			message := model_builders.NewChatCompletionMessageBuilder().Build()

			initSuite(suite, llmConnection, []*models.ChatCompletionMessage{message})

			result, err := suite.sut.ChatCompletion(llmConnection, nil, []*models.ChatCompletionMessage{message}, nil)

			assert.Nil(t, err)
			assert.NotNil(t, result)
		})
	})

	t.Run("Given llm connection is Mistral", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingChatCompletingTestingSuite, llmConnection *models.LLMConnection, messages []*models.ChatCompletionMessage, apiKey string) {
			suite.mistralSettings.BaseUrl = "http://mistral.api"

			result := "completion-result"

			suite.mockOpenAICompletionsStandardChatCompletionRequester.
				EXPECT().
				ChatCompletion(
					suite.mistralSettings.GetBaseUrl(),
					apiKey,
					llmConnection,
					gomock.Any(),
					nil,
				).
				DoAndReturn(func(_, _ string, _ *models.LLMConnection, externalMessages []external_models.ExternalOpenAICompletionsStandardRequestMessage, _ []*models.MCPConnection) (*string, error) {
					assert.Equal(t, len(messages), len(externalMessages))
					for i, message := range messages {
						assert.Equal(t, message.Role, externalMessages[i].Role)
						assert.Equal(t, message.Content, *externalMessages[i].Content)
					}

					return &result, nil
				})
		}

		t.Run("Should use OpenAI completions standard", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingChatCompletionBeforeEach(t)

			llmConnection := model_builders.NewLLMConnectionBuilder().
				WithPlatform(enums.LLMConnectionPlatformMistral).
				WithModel("mistral-7b").
				Build()
			message := model_builders.NewChatCompletionMessageBuilder().Build()
			apiKey := "api-key"

			initSuite(suite, llmConnection, []*models.ChatCompletionMessage{message}, apiKey)

			result, err := suite.sut.ChatCompletion(llmConnection, nil, []*models.ChatCompletionMessage{message}, &apiKey)

			assert.Nil(t, err)
			assert.NotNil(t, result)
		})
	})

	t.Run("Given llm connection is DeepSeek", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingChatCompletingTestingSuite, llmConnection *models.LLMConnection, messages []*models.ChatCompletionMessage, apiKey string) {
			suite.deepSeekSettings.BaseUrl = "http://deepseek.api"

			result := "completion-result"

			suite.mockOpenAICompletionsStandardChatCompletionRequester.
				EXPECT().
				ChatCompletion(
					suite.deepSeekSettings.GetBaseUrl(),
					apiKey,
					llmConnection,
					gomock.Any(),
					nil,
				).
				DoAndReturn(func(_, _ string, _ *models.LLMConnection, externalMessages []external_models.ExternalOpenAICompletionsStandardRequestMessage, _ []*models.MCPConnection) (*string, error) {
					assert.Equal(t, len(messages), len(externalMessages))
					for i, message := range messages {
						assert.Equal(t, message.Role, externalMessages[i].Role)
						assert.Equal(t, message.Content, *externalMessages[i].Content)
					}

					return &result, nil
				})
		}

		t.Run("Should use OpenAI completions standard", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingChatCompletionBeforeEach(t)

			llmConnection := model_builders.NewLLMConnectionBuilder().
				WithPlatform(enums.LLMConnectionPlatformDeepseek).
				WithModel("deepseek-7b").
				Build()
			message := model_builders.NewChatCompletionMessageBuilder().Build()
			apiKey := "api-key"

			initSuite(suite, llmConnection, []*models.ChatCompletionMessage{message}, apiKey)

			result, err := suite.sut.ChatCompletion(llmConnection, nil, []*models.ChatCompletionMessage{message}, &apiKey)

			assert.Nil(t, err)
			assert.NotNil(t, result)
		})
	})

	t.Run("Given llm connection is X", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingChatCompletingTestingSuite, llmConnection *models.LLMConnection, messages []*models.ChatCompletionMessage, apiKey string) {
			suite.xSettings.BaseUrl = "http://x.api"

			result := "completion-result"

			suite.mockOpenAICompletionsStandardChatCompletionRequester.
				EXPECT().
				ChatCompletion(
					suite.xSettings.GetBaseUrl(),
					apiKey,
					llmConnection,
					gomock.Any(),
					nil,
				).
				DoAndReturn(func(_, _ string, _ *models.LLMConnection, externalMessages []external_models.ExternalOpenAICompletionsStandardRequestMessage, _ []*models.MCPConnection) (*string, error) {
					assert.Equal(t, len(messages), len(externalMessages))
					for i, message := range messages {
						assert.Equal(t, message.Role, externalMessages[i].Role)
						assert.Equal(t, message.Content, *externalMessages[i].Content)
					}

					return &result, nil
				})
		}

		t.Run("Should use OpenAI completions standard", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingChatCompletionBeforeEach(t)

			llmConnection := model_builders.NewLLMConnectionBuilder().
				WithPlatform(enums.LLMConnectionPlatformX).
				WithModel("grok").
				Build()
			message := model_builders.NewChatCompletionMessageBuilder().Build()
			apiKey := "api-key"

			initSuite(suite, llmConnection, []*models.ChatCompletionMessage{message}, apiKey)

			result, err := suite.sut.ChatCompletion(llmConnection, nil, []*models.ChatCompletionMessage{message}, &apiKey)

			assert.Nil(t, err)
			assert.NotNil(t, result)
		})
	})

	t.Run("Given llm connection is Anthropic", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingChatCompletingTestingSuite, llmConnection *models.LLMConnection, messages []*models.ChatCompletionMessage, apiKey string) {
			suite.anthropicSettings.BaseUrl = "http://anthropic.api"

			result := "completion-result"

			suite.mockAnthropicMessagesStandardChatCompletionRequester.
				EXPECT().
				ChatCompletion(
					suite.anthropicSettings.GetBaseUrl(),
					apiKey,
					suite.anthropicSettings.HasDirectAccessToMcp(),
					llmConnection,
					gomock.Any(),
					gomock.Any(),
					nil,
				).
				DoAndReturn(func(_, _ string, _ bool, _ *models.LLMConnection, externalMessages []external_models.ExternalAnthropicMessagesStandardMessage, _ []external_models.ExternalAnthropicMessagesStandardRequestSystem, _ []*models.MCPConnection) (*string, error) {
					assert.Equal(t, len(messages), len(externalMessages))
					for i, message := range messages {
						assert.Equal(t, message.Role, externalMessages[i].Role)
						assert.Equal(t, message.Content, *externalMessages[i].Content[0].Text)
					}

					return &result, nil
				})
		}

		t.Run("Should use Anthropic messages standard", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingChatCompletionBeforeEach(t)

			llmConnection := model_builders.NewLLMConnectionBuilder().
				WithPlatform(enums.LLMConnectionPlatformAnthropic).
				WithModel("claude").
				Build()
			message := model_builders.NewChatCompletionMessageBuilder().Build()
			apiKey := "api-key"

			initSuite(suite, llmConnection, []*models.ChatCompletionMessage{message}, apiKey)

			result, err := suite.sut.ChatCompletion(llmConnection, nil, []*models.ChatCompletionMessage{message}, &apiKey)

			assert.Nil(t, err)
			assert.NotNil(t, result)
		})
	})

	t.Run("Given llm connection is Amazon Bedrock", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRequestingChatCompletingTestingSuite, llmConnection *models.LLMConnection, messages []*models.ChatCompletionMessage, apiKey string) {
			result := "completion-result"

			suite.mockAmazonBedrockConverseStandardChatCompletionRequester.
				EXPECT().
				ChatCompletion(
					*llmConnection.BaseUrl,
					apiKey,
					llmConnection,
					gomock.Any(),
					gomock.Any(),
					nil,
				).
				DoAndReturn(func(_, _ string, _ *models.LLMConnection, externalMessages []external_models.ExternalAmazonBedrockConverseStandardMessage, externalSystemMessages []external_models.ExternalAmazonBedrockConverseStandardRequestSystem, _ []*models.MCPConnection) (*string, error) {
					assert.Equal(t, len(messages), len(externalMessages))
					for i, message := range messages {
						assert.Equal(t, message.Role, externalMessages[i].Role)
						assert.Equal(t, message.Content, *externalMessages[i].Content[0].Text)
					}

					return &result, nil
				})
		}

		t.Run("Should use Anthropic messages standard", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingChatCompletionBeforeEach(t)

			llmConnection := model_builders.NewLLMConnectionBuilder().
				WithPlatform(enums.LLMConnectionPlatformAmazonBedrock).
				WithModel("gpt-4").
				Build()
			message := model_builders.NewChatCompletionMessageBuilder().Build()
			apiKey := "api-key"

			initSuite(suite, llmConnection, []*models.ChatCompletionMessage{message}, apiKey)

			result, err := suite.sut.ChatCompletion(llmConnection, nil, []*models.ChatCompletionMessage{message}, &apiKey)

			assert.Nil(t, err)
			assert.NotNil(t, result)
		})
	})

	t.Run("Given llm connection is not supported", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequestingChatCompletionBeforeEach(t)

			llmConnection := model_builders.NewLLMConnectionBuilder().
				WithPlatform("unsupported").
				WithModel("claude").
				Build()
			message := model_builders.NewChatCompletionMessageBuilder().Build()
			apiKey := "api-key"

			result, err := suite.sut.ChatCompletion(llmConnection, nil, []*models.ChatCompletionMessage{message}, &apiKey)

			assert.NotNil(t, err)
			assert.Nil(t, result)
		})
	})
}
