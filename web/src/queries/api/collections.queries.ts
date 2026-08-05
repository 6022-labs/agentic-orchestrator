import axios from "axios"

import type { AgentCollection } from "@/models/agent-collection"

// GET /agent-collections — collections the orchestrator manages on-chain.
export const listAgentCollections = async () => {
  return await axios.get<AgentCollection[]>("/agent-collections")
}
