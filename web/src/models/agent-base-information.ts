import { AgentRole } from "./enums/agent-role";
import { AgentBaseLLMSettings } from "./agent-base-llm-settings";

export type AgentBaseInformation = AgentBaseLLMSettings & {
  name: string;
  role: AgentRole;
};
