// Live 6022 well-known card the gateway relays from the agent's node.
export type AgentCard = {
  protocol: string
  protocolVersion: string
  name: string
  ensDomain?: string
  endpoints: {
    responses?: string
    chatCompletions?: string
  }
  paymentMethods?: string[]
}
