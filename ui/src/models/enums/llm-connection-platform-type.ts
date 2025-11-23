import { AgentRole } from "./agent-role";

export enum LLMConnectionPlatformType {
  None = "none", // Used for default value
  Cloud = "cloud",
  SelfHost = "self-hosted",
}

export const LLMConnectionPlatformTypeDisplayNames: Record<
  LLMConnectionPlatformType,
  string
> = {
  [LLMConnectionPlatformType.None]: "",
  [LLMConnectionPlatformType.Cloud]: "Cloud",
  [LLMConnectionPlatformType.SelfHost]: "Self Hosted",
};

export const getPossibleLLMConnectionPlatformTypeByAgentRole = (
  agentRole: AgentRole
) => {
  switch (agentRole) {
    case AgentRole.Human:
      return [LLMConnectionPlatformType.None];
    default:
      return [
        LLMConnectionPlatformType.Cloud,
        LLMConnectionPlatformType.SelfHost,
      ];
  }
};
