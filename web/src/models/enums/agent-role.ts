export enum AgentRole {
  None = "none", // Used for default value
  Clone = "clone",
  Human = "human",
  Expert = "expert",
  Facilitator = "facilitator",
}

export const AgentRoleDisplayNames: Record<AgentRole, string> = {
  [AgentRole.None]: "",
  [AgentRole.Clone]: "Digital twin",
  [AgentRole.Human]: "Human",
  [AgentRole.Expert]: "Expert",
  [AgentRole.Facilitator]: "Facilitator",
};

export const getAllAgentRoles = (): AgentRole[] => {
  return [
    AgentRole.Clone,
    AgentRole.Human,
    AgentRole.Expert,
    AgentRole.Facilitator,
  ];
};

export const getAllAgentRolesExceptClone = (): AgentRole[] => {
  return getAllAgentRoles().filter((role) => role !== AgentRole.Clone);
};
