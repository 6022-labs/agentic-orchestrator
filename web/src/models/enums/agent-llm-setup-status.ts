export enum AgentLLMSetupStatus {
  Configured = "configured",
  Established = "established",
}

export const AgentLLMSetupStatusDisplayNames: Record<
  AgentLLMSetupStatus,
  string
> = {
  [AgentLLMSetupStatus.Configured]: "Configured",
  [AgentLLMSetupStatus.Established]: "Established",
};
