import { LLMConnectionPlatform } from "./enums/llm-connection-platform";

export type AgentBaseLLMSettings = {
  baseUrl?: string;
  model: string;
  platform: LLMConnectionPlatform;
};
