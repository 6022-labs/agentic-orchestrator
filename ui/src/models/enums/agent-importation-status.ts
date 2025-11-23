export enum AgentImportationStatus {
  FirstStep = "first-step",
  SecondStep = "second-step",
  ThirdStep = "third-step",
  Finished = "finished",
}

export const AgentImportationStatusDisplayNames: Record<
  AgentImportationStatus,
  string
> = {
  [AgentImportationStatus.FirstStep]: "Pending",
  [AgentImportationStatus.SecondStep]: "Pending",
  [AgentImportationStatus.ThirdStep]: "Pending",
  [AgentImportationStatus.Finished]: "Finished",
};
