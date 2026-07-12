import axios from "axios"

import type { Agent, ListAgentsResponse } from "@/models/agent"
import type { AgentCard } from "@/models/agent-card"

export type ListAgentsParams = {
  page?: number
  elementsPerPage?: number
  sortBy?: string
  sortDirection?: string
}

// GET /agents — paginated agent directory.
export const listAgents = async (params: ListAgentsParams = {}) => {
  return await axios.get<ListAgentsResponse>("/agents", { params })
}

// GET /agent/{collectionAddress}/{collectionTokenId} — one agent's identity.
export const getAgent = async (
  collectionAddress: string,
  collectionTokenId: string
) => {
  return await axios.get<Agent>(
    `/agent/${collectionAddress}/${collectionTokenId}`
  )
}

// GET /agent/{collectionAddress}/{collectionTokenId}/card — the live 6022 card
// relayed from the agent's node; fails when the node is unreachable.
export const getAgentCard = async (
  collectionAddress: string,
  collectionTokenId: string
) => {
  return await axios.get<AgentCard>(
    `/agent/${collectionAddress}/${collectionTokenId}/card`
  )
}
