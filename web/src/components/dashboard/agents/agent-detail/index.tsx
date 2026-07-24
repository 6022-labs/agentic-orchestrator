import React from "react"
import { toast } from "sonner"
import { Link, useParams } from "react-router-dom"
import { ArrowLeft, Bot, ExternalLink } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import EmptyState from "@/components/dashboard/empty-state"
import AgentAvatar from "@/components/dashboard/agents/agent-avatar"
import AgentCardPanel from "@/components/dashboard/agents/agent-detail/agent-card-panel"
import PromptTemplatesPanel from "@/components/dashboard/agents/agent-detail/prompt-templates-panel"
import TransferAgentDialog from "@/components/dashboard/agents/dialogs/transfer-agent-dialog"
import { openAgentNode } from "@/lib/agent-node"
import { extractErrorMessage } from "@/lib/api-errors"
import { truncateEthAddress } from "@/lib/utils"
import type { Agent } from "@/models/agent"
import { getAgent } from "@/queries/api/agents.queries"

function AgentField({ label, value }: { label: string; value?: string }) {
  return (
    <div className="flex flex-col">
      <span className="text-xs text-muted-foreground">{label}</span>
      <span className="truncate text-sm font-medium">{value || "—"}</span>
    </div>
  )
}

export default function AgentDetailSection() {
  const { collectionAddress, collectionTokenId } = useParams<{
    collectionAddress: string
    collectionTokenId: string
  }>()

  const [loading, setLoading] = React.useState(true)
  const [agent, setAgent] = React.useState<Agent | null>(null)

  const load = React.useCallback(async () => {
    if (!collectionAddress || !collectionTokenId) return

    setLoading(true)
    try {
      const { data } = await getAgent(collectionAddress, collectionTokenId)
      setAgent(data)
    } catch (error) {
      toast.error(extractErrorMessage(error, "Failed to load the agent"))
      setAgent(null)
    } finally {
      setLoading(false)
    }
  }, [collectionAddress, collectionTokenId])

  React.useEffect(() => {
    void load()
  }, [load])

  if (!collectionAddress || !collectionTokenId) {
    return (
      <EmptyState
        icon={Bot}
        title="Agent not found"
        description="The agent route is missing its collection or token id."
      />
    )
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <Button variant="ghost" size="sm" asChild>
          <Link to="/dashboard/agents">
            <ArrowLeft />
            Back to agents
          </Link>
        </Button>
        {agent ? (
          <div className="flex items-center gap-2">
            {agent.ensDomain ? (
              <Button
                variant="ghost"
                size="icon"
                className="size-6"
                title="Open agent node"
                aria-label="Open agent node"
                onClick={() =>
                  openAgentNode(collectionAddress, collectionTokenId)
                }
              >
                <ExternalLink className="size-3.5" />
              </Button>
            ) : null}
            <TransferAgentDialog
              collectionAddress={collectionAddress}
              collectionTokenId={collectionTokenId}
              agentName={agent.name}
              onTransferred={() => void load()}
            />
          </div>
        ) : null}
      </div>

      {loading ? (
        <Card>
          <CardContent className="flex items-center gap-4">
            <Skeleton className="size-16 rounded-full" />
            <div className="flex-1 space-y-2">
              <Skeleton className="h-6 w-48" />
              <Skeleton className="h-4 w-64" />
            </div>
          </CardContent>
        </Card>
      ) : !agent ? (
        <EmptyState
          icon={Bot}
          title="Agent not found"
          description="The orchestrator has no agent at this address and token id."
        />
      ) : (
        <>
          <Card>
            <CardContent className="flex flex-col gap-6">
              <div className="flex items-center gap-4">
                <AgentAvatar agent={agent} className="size-16" />
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <h2 className="truncate text-xl font-semibold">
                      {agent.name}
                    </h2>
                    <Badge variant="secondary">
                      {agent.attributes["role"]}
                    </Badge>
                  </div>
                  <p className="truncate text-sm text-muted-foreground">
                    {agent.ensDomain || "No ENS domain"}
                  </p>
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
                <AgentField
                  label="Owner"
                  value={truncateEthAddress(agent.owner)}
                />
                <AgentField
                  label="Wallet"
                  value={truncateEthAddress(agent.addresses?.["evm"])}
                />
                <AgentField
                  label="Collection"
                  value={truncateEthAddress(agent.collectionAddress)}
                />
                <AgentField
                  label="Token ID"
                  value={agent.collectionTokenId?.toString()}
                />
              </div>

              {Object.keys(agent.attributes ?? {}).length > 0 ? (
                <div className="flex flex-col gap-2">
                  <span className="text-xs text-muted-foreground">
                    Attributes
                  </span>
                  <div className="flex flex-wrap gap-2">
                    {Object.entries(agent.attributes).map(([key, value]) => (
                      <Badge key={key} variant="outline">
                        {key}: {value}
                      </Badge>
                    ))}
                  </div>
                </div>
              ) : null}
            </CardContent>
          </Card>

          <div className="grid gap-6 lg:grid-cols-2">
            <AgentCardPanel
              collectionAddress={collectionAddress}
              collectionTokenId={collectionTokenId}
            />
            <PromptTemplatesPanel
              collectionAddress={collectionAddress}
              collectionTokenId={collectionTokenId}
              role={agent.attributes["role"] ?? ""}
            />
          </div>
        </>
      )}
    </div>
  )
}
