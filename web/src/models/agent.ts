import { AgentRole } from "./enums/agent-role";
import { AgentLLMSetupStatus } from "./enums/agent-llm-setup-status";
import { AgentDatabaseSetupStatus } from "./enums/agent-database-setup-status";

export type Agent = {
  id: string;
  role: AgentRole;
  owner: string;
  creator: string;
  images: Record<string, string>;
  attributes: Record<string, string>;
  cloneOf: number | null;
  isInternal: boolean;
  walletAddress: string;
  name: string | undefined;
  collectionAddress: string;
  isSetupCompleted: boolean | undefined;
  proposalId: number | undefined;
  collectionTokenId: number | undefined;
  llmSetupStatus: AgentLLMSetupStatus | undefined;
  databaseSetupStatus: AgentDatabaseSetupStatus | undefined;
};
