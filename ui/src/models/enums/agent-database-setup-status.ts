export enum AgentDatabaseSetupStatus {
  Configured = "configured",
  Established = "established",
  Establishing = "establishing",
}

export const AgentDatabaseSetupStatusDisplayNames: Record<
  AgentDatabaseSetupStatus,
  string
> = {
  [AgentDatabaseSetupStatus.Configured]: "Configured",
  [AgentDatabaseSetupStatus.Established]: "Established",
  [AgentDatabaseSetupStatus.Establishing]: "Establishing",
};
