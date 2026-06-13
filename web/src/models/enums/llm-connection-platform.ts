import { LLMConnectionPlatformType } from "./llm-connection-platform-type";

export enum LLMConnectionPlatform {
  X = "x",
  None = "none", // Used for the default value
  OpenAI = "open-ai",
  Ollama = "ollama",
  Mistral = "mistral",
  Deepseek = "deepseek",
  LibertAI = "libert-ai",
  Anthropic = "anthropic",
  AmazonBedrock = "amazon-bedrock",
}

export const LLMConnectionPlatformDisplayNames: Record<
  LLMConnectionPlatform,
  string
> = {
  [LLMConnectionPlatform.X]: "X",
  [LLMConnectionPlatform.None]: "",
  [LLMConnectionPlatform.OpenAI]: "OpenAI",
  [LLMConnectionPlatform.Ollama]: "Ollama",
  [LLMConnectionPlatform.Mistral]: "Mistral",
  [LLMConnectionPlatform.LibertAI]: "LibertAI",
  [LLMConnectionPlatform.Deepseek]: "Deepseek",
  [LLMConnectionPlatform.Anthropic]: "Anthropic",
  [LLMConnectionPlatform.AmazonBedrock]: "Amazon Bedrock",
};

export function getAllSelfHostedLLMConnectionPlatforms(): LLMConnectionPlatform[] {
  return [LLMConnectionPlatform.Ollama];
}

export function getAllCloudLLMConnectionPlatforms(): LLMConnectionPlatform[] {
  return [
    LLMConnectionPlatform.X,
    LLMConnectionPlatform.OpenAI,
    LLMConnectionPlatform.Mistral,
    LLMConnectionPlatform.Deepseek,
    LLMConnectionPlatform.LibertAI,
    LLMConnectionPlatform.Anthropic,
    LLMConnectionPlatform.AmazonBedrock,
  ];
}

export function getLLMConnectionPlatformType(
  platform: LLMConnectionPlatform
): LLMConnectionPlatformType {
  if (getAllSelfHostedLLMConnectionPlatforms().includes(platform)) {
    return LLMConnectionPlatformType.SelfHost;
  }

  return LLMConnectionPlatformType.Cloud;
}
