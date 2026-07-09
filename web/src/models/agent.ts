import { AgentRole } from "./enums/agent-role";

export type Agent = {
  id: string;
  role: AgentRole;
  owner: string;
  creator: string;
  images: Record<string, string>;
  attributes: Record<string, string>;
  cloneOf: number | null;
  walletAddress: string;
  name: string | undefined;
  collectionAddress: string;
  proposalId: number | undefined;
  collectionTokenId: number | undefined;
};
