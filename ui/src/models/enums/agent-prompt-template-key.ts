export enum AgentPromptTemplateKey {
  CloneSystem = "clone-system",
  CloneDefault = "clone-default",
  FacilitatorSystem = "facilitator-system",
  FacilitatorDefault = "facilitator-default",
  FacilitatorConclude = "facilitator-conclude",
  ExpertSystem = "expert-system",
  ExpertContextFound = "expert-context-found",
  ExpertContextNotFound = "expert-context-not-found",
}

export const AgentPromptTemplateKeyDisplayNames: Record<
  AgentPromptTemplateKey,
  string
> = {
  [AgentPromptTemplateKey.CloneSystem]: "System",
  [AgentPromptTemplateKey.CloneDefault]: "Default",
  [AgentPromptTemplateKey.FacilitatorSystem]: "System",
  [AgentPromptTemplateKey.FacilitatorDefault]: "Default",
  [AgentPromptTemplateKey.FacilitatorConclude]: "Conclude",
  [AgentPromptTemplateKey.ExpertSystem]: "System",
  [AgentPromptTemplateKey.ExpertContextFound]: "Context Found",
  [AgentPromptTemplateKey.ExpertContextNotFound]: "Context Not Found",
};

export const FacilitatorAgentPromptTemplateKeys = [
  AgentPromptTemplateKey.FacilitatorSystem,
  AgentPromptTemplateKey.FacilitatorDefault,
  AgentPromptTemplateKey.FacilitatorConclude,
];

export const ExpertAgentPromptTemplateKeys = [
  AgentPromptTemplateKey.ExpertSystem,
  AgentPromptTemplateKey.ExpertContextFound,
  AgentPromptTemplateKey.ExpertContextNotFound,
];

export const CloneAgentPromptTemplateKeys = [
  AgentPromptTemplateKey.CloneSystem,
  AgentPromptTemplateKey.CloneDefault,
];

export const GetAgentPromptTemplateKeysFromRole = (
  role: string
): AgentPromptTemplateKey[] => {
  if (role === "facilitator") {
    return FacilitatorAgentPromptTemplateKeys;
  }

  if (role === "clone") {
    return CloneAgentPromptTemplateKeys;
  }

  return ExpertAgentPromptTemplateKeys;
};
