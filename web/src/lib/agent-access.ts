import type { Agent } from "@/models/agent"

// Mirrors the smart-contract / agent-node rule: the owner and the agent's own
// wallets may act on an agent.
export const canManageAgent = (agent: Agent, address?: string) => {
  if (!address) return false
  const wallet = address.toLowerCase()
  if (agent.owner?.toLowerCase() === wallet) return true
  return Object.values(agent.addresses ?? {}).some(
    (agentAddress) => agentAddress.toLowerCase() === wallet
  )
}
