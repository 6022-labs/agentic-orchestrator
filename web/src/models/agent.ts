import type { AgentRole } from "@/models/enums/agent-role"

export type Agent = {
  id: string
  name: string
  role: AgentRole
  owner?: string
  creator: string
  // Image key -> IPFS CID (e.g. "default").
  images: Record<string, string>
  cloneOf?: string
  attributes: Record<string, string>
  ensDomain: string
  // Address type (e.g. "evm", "btc") -> address value.
  addresses: Record<string, string>
  proposalId?: number
  collectionTokenId?: number
  collectionAddress: string
}

export type ListAgentsResponse = {
  page: number
  elementsPerPage: number
  totalElements: number
  agents: Agent[]
}
