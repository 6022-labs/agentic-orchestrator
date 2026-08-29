import React from "react"
import { toast } from "sonner"
import { Link, useParams } from "react-router-dom"
import { ArrowLeft, Bot, ExternalLink } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button, buttonVariants } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import EmptyState from "@/components/dashboard/empty-state"
import AgentCardPanel from "@/components/dashboard/agents/agent-detail/agent-card-panel"
import PromptTemplatesPanel from "@/components/dashboard/agents/agent-detail/prompt-templates-panel"
import TransferAgentDialog from "@/components/dashboard/agents/dialogs/transfer-agent-dialog"
import { canManageAgent } from "@/lib/agent-access"
import { openAgentNode } from "@/lib/agent-node"
import { extractErrorMessage } from "@/lib/api-errors"
import { useWallet } from "@/contexts/wallet-provider-context"
import IpfsImage from "@/components/ipfs-image"
import { chainIdFromEnsDomain } from "@/lib/ens-label"
import {
  getEnsAppUrl,
  getOpenSeaUrl,
  getRaribleUrl,
  getScannerAddressUrl,
} from "@/configs/networks"
import type { Agent } from "@/models/agent"
import { getAgent } from "@/queries/api/agents.queries"

function Row({
  label,
  value,
  mono,
  href,
}: {
  label: string
  value: string
  mono?: boolean
  href?: string | null
}) {
  return (
    <div className="grid grid-cols-[7.5rem_1fr] items-baseline gap-4 px-4 py-3.5 sm:grid-cols-[9rem_1fr]">
      <dt className="text-sm text-muted-foreground">{label}</dt>
      <dd className="min-w-0">
        {href ? (
          <a
            href={href}
            target="_blank"
            rel="noreferrer"
            className={`flex items-center gap-1 text-sm hover:underline ${mono ? "font-mono" : ""}`}
          >
            <span className="truncate">{value}</span>
            <ExternalLink className="size-3.5 shrink-0" />
          </a>
        ) : (
          <span className={`block truncate text-sm ${mono ? "font-mono" : ""}`}>
            {value}
          </span>
        )}
      </dd>
    </div>
  )
}

export default function AgentDetailSection() {
  const { collectionAddress, collectionTokenId } = useParams<{
    collectionAddress: string
    collectionTokenId: string
  }>()

  const { address } = useWallet()
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

  const chainId = agent ? chainIdFromEnsDomain(agent.ensDomain) : ""
  const imageValue = agent
    ? (agent.images?.default ?? Object.values(agent.images ?? {})[0])
    : null
  const tokenId = agent?.collectionTokenId ?? null
  const openSeaUrl = agent
    ? getOpenSeaUrl(chainId, agent.collectionAddress, tokenId)
    : null
  const raribleUrl = agent
    ? getRaribleUrl(chainId, agent.collectionAddress, tokenId)
    : null

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <Button variant="ghost" size="sm" asChild>
          <Link to="/dashboard/agents">
            <ArrowLeft />
            Back to agents
          </Link>
        </Button>
        {agent && canManageAgent(agent, address) ? (
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
          <CardContent className="flex flex-col gap-6 sm:flex-row">
            <Skeleton className="aspect-square w-full shrink-0 rounded-lg sm:w-64" />
            <div className="flex flex-1 flex-col justify-center gap-4">
              <Skeleton className="h-8 w-48" />
              <div className="space-y-2">
                <Skeleton className="h-4 w-full" />
                <Skeleton className="h-4 w-3/4" />
                <Skeleton className="h-4 w-2/3" />
              </div>
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
              <div className="flex flex-col gap-6 sm:flex-row sm:items-stretch">
                {imageValue || openSeaUrl || raribleUrl ? (
                  <div className="flex w-full shrink-0 flex-col gap-2 sm:w-64">
                    <IpfsImage
                      value={imageValue}
                      alt={agent.name}
                      className="w-full rounded-lg border"
                      fallback={
                        <div className="flex aspect-square w-full items-center justify-center rounded-lg border bg-muted text-muted-foreground">
                          <Bot className="size-1/3" />
                        </div>
                      }
                    />
                    {openSeaUrl || raribleUrl ? (
                      <div className="grid grid-cols-2 gap-2">
                        {openSeaUrl ? (
                          <a
                            href={openSeaUrl}
                            target="_blank"
                            rel="noreferrer"
                            className={buttonVariants({
                              variant: "outline",
                              size: "sm",
                              className: "w-full",
                            })}
                          >
                            OpenSea
                            <ExternalLink />
                          </a>
                        ) : null}
                        {raribleUrl ? (
                          <a
                            href={raribleUrl}
                            target="_blank"
                            rel="noreferrer"
                            className={buttonVariants({
                              variant: "outline",
                              size: "sm",
                              className: "w-full",
                            })}
                          >
                            Rarible
                            <ExternalLink />
                          </a>
                        ) : null}
                      </div>
                    ) : null}
                  </div>
                ) : null}
                <div className="flex flex-1 flex-col justify-center gap-6">
                  <div className="flex flex-wrap items-center gap-x-3 gap-y-2">
                    <h2 className="text-2xl font-semibold tracking-tight sm:text-3xl">
                      {agent.name}
                    </h2>
                    {agent.attributes?.["role"] ? (
                      <Badge variant="secondary" className="capitalize">
                        {agent.attributes["role"]}
                      </Badge>
                    ) : null}
                  </div>
                  <dl className="divide-y rounded-lg border">
                    <Row
                      label="Collection"
                      value={agent.collectionAddress}
                      mono
                      href={getScannerAddressUrl(
                        chainId,
                        agent.collectionAddress
                      )}
                    />
                    <Row
                      label="Token id"
                      value={tokenId != null ? `#${tokenId}` : "—"}
                    />
                    {agent.addresses?.["evm"] ? (
                      <Row
                        label="Agent wallet"
                        value={agent.addresses["evm"]}
                        mono
                        href={getScannerAddressUrl(
                          chainId,
                          agent.addresses["evm"]
                        )}
                      />
                    ) : null}
                    {agent.owner ? (
                      <Row
                        label="Owner"
                        value={agent.owner}
                        mono
                        href={getScannerAddressUrl(chainId, agent.owner)}
                      />
                    ) : null}
                    {agent.ensDomain ? (
                      <Row
                        label="ENS domain"
                        value={agent.ensDomain}
                        href={getEnsAppUrl(chainId, agent.ensDomain)}
                      />
                    ) : null}
                  </dl>
                </div>
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
