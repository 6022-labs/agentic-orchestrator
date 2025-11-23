export enum AgentContextMessageRole {
  User = "user",
  System = "system",
  Assistant = "assistant",
}

export const AgentContextMessageRoleDisplayNames: Record<
  AgentContextMessageRole,
  string
> = {
  [AgentContextMessageRole.User]: "User",
  [AgentContextMessageRole.System]: "System",
  [AgentContextMessageRole.Assistant]: "Assistant",
};
