package enums

import "strings"

type LLMConnectionPlatform string
type LLMConnectionPlatformList []LLMConnectionPlatform

const (
	LLMConnectionPlatformX             LLMConnectionPlatform = "x" // Twitter
	LLMConnectionPlatformMistral       LLMConnectionPlatform = "mistral"
	LLMConnectionPlatformOpenAI        LLMConnectionPlatform = "open-ai"
	LLMConnectionPlatformLibertAI      LLMConnectionPlatform = "libert-ai"
	LLMConnectionPlatformOllama        LLMConnectionPlatform = "ollama"
	LLMConnectionPlatformDeepseek      LLMConnectionPlatform = "deepseek"
	LLMConnectionPlatformAnthropic     LLMConnectionPlatform = "anthropic"
	LLMConnectionPlatformAmazonBedrock LLMConnectionPlatform = "amazon-bedrock"
)

func GetAllLLMConnectionPlatforms() LLMConnectionPlatformList {
	return LLMConnectionPlatformList{
		LLMConnectionPlatformX,
		LLMConnectionPlatformOpenAI,
		LLMConnectionPlatformOllama,
		LLMConnectionPlatformMistral,
		LLMConnectionPlatformDeepseek,
		LLMConnectionPlatformLibertAI,
		LLMConnectionPlatformAnthropic,
		LLMConnectionPlatformAmazonBedrock,
	}
}

func GetAllLocalLLMConnectionPlatforms() LLMConnectionPlatformList {
	return LLMConnectionPlatformList{
		LLMConnectionPlatformOllama,
	}
}

func GetAllCloudLLMConnectionPlatforms() LLMConnectionPlatformList {
	return LLMConnectionPlatformList{
		LLMConnectionPlatformX,
		LLMConnectionPlatformOpenAI,
		LLMConnectionPlatformMistral,
		LLMConnectionPlatformDeepseek,
		LLMConnectionPlatformLibertAI,
		LLMConnectionPlatformAnthropic,
		LLMConnectionPlatformAmazonBedrock,
	}
}

func (k LLMConnectionPlatform) String() string {
	return string(k)
}

func (k LLMConnectionPlatformList) String() string {
	var platforms []string
	for _, v := range k {
		platforms = append(platforms, string(v))
	}
	return strings.Join(platforms, ", ")
}

func (k LLMConnectionPlatformList) Contains(platform LLMConnectionPlatform) bool {
	for _, v := range k {
		if v == platform {
			return true
		}
	}
	return false
}
