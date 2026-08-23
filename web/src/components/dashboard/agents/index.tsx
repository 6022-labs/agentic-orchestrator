import React from "react"
import { useNavigate } from "react-router-dom"
import { Bot, ChevronLeft, ChevronRight, ExternalLink } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import EmptyState from "@/components/dashboard/empty-state"
import CreateCollectionDialog from "@/components/dashboard/agents/dialogs/create-collection-dialog"
import TransferAgentDialog from "@/components/dashboard/agents/dialogs/transfer-agent-dialog"
import AgentAvatar from "@/components/dashboard/agents/agent-avatar"
import { useAsyncData } from "@/hooks/use-async-data"
import { truncateEthAddress } from "@/lib/utils"
import { canManageAgent } from "@/lib/agent-access"
import { openAgentNode } from "@/lib/agent-node"
import { listAgents } from "@/queries/api/agents.queries"
import { listAgentCollections } from "@/queries/api/collections.queries"
import { useWallet } from "@/contexts/wallet-provider-context"

const ELEMENTS_PER_PAGE = 20

type AgentFilter = "all" | "minted" | "proposals"

const FILTERS: {
  value: AgentFilter
  label: string
  isMinted?: boolean
  noun: string
  emptyDescription: string
}[] = [
  {
    value: "all",
    label: "All",
    noun: "agents",
    emptyDescription: "Registered agents show up here once they are proposed.",
  },
  {
    value: "minted",
    label: "Minted",
    isMinted: true,
    noun: "minted agents",
    emptyDescription: "Registered agents show up here once they are minted.",
  },
  {
    value: "proposals",
    label: "Proposals",
    isMinted: false,
    noun: "proposals",
    emptyDescription: "No agents are waiting to be minted.",
  },
]

export default function AgentsSection() {
  const navigate = useNavigate()
  const { address } = useWallet()
  const [page, setPage] = React.useState(1)
  const [filter, setFilter] = React.useState<AgentFilter>("minted")

  const activeFilter = FILTERS.find((f) => f.value === filter) ?? FILTERS[1]

  const agentsSource = useAsyncData(
    React.useCallback(
      () =>
        listAgents({
          page,
          elementsPerPage: ELEMENTS_PER_PAGE,
          isMinted: activeFilter.isMinted,
        }),
      [page, activeFilter]
    ),
    "Failed to load agents"
  )
  // Collections only name the agent rows — a failing lookup must not blank the table.
  const collectionsSource = useAsyncData(
    React.useCallback(() => listAgentCollections(), []),
    "Failed to load collections"
  )

  const loading = agentsSource.loading
  const agents = agentsSource.data?.data.agents ?? []
  const totalElements = agentsSource.data?.data.totalElements ?? 0
  const collections = collectionsSource.data?.data ?? []

  const collectionName = (address: string) =>
    collections.find(
      (collection) => collection.address.toLowerCase() === address.toLowerCase()
    )?.name ?? truncateEthAddress(address)

  const totalPages = Math.max(Math.ceil(totalElements / ELEMENTS_PER_PAGE), 1)

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between gap-2">
          <div>
            <CardTitle>Agents</CardTitle>
            <CardDescription>
              Every agent identity the orchestrator knows about.
            </CardDescription>
          </div>
          <CreateCollectionDialog
            onCreated={() => {
              agentsSource.reload()
              collectionsSource.reload()
            }}
          />
        </div>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        <div className="flex gap-2">
          {FILTERS.map((option) => (
            <Button
              key={option.value}
              variant={option.value === filter ? "default" : "outline"}
              size="sm"
              onClick={() => {
                setFilter(option.value)
                setPage(1)
              }}
            >
              {option.label}
            </Button>
          ))}
        </div>

        {loading ? (
          <div className="space-y-2">
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
          </div>
        ) : agents.length === 0 ? (
          <EmptyState
            icon={Bot}
            title="No agents yet"
            description={activeFilter.emptyDescription}
          />
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-12"></TableHead>
                <TableHead>Name</TableHead>
                <TableHead>ENS domain</TableHead>
                <TableHead>Role</TableHead>
                <TableHead>Owner</TableHead>
                <TableHead>Collection</TableHead>
                <TableHead className="w-20"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {agents.map((agent) => {
                const isProposal = agent.collectionTokenId == null

                return (
                  <TableRow
                    key={agent.id}
                    className={isProposal ? undefined : "cursor-pointer"}
                    onClick={(event) => {
                      if (isProposal) return
                      if (!event.currentTarget.contains(event.target as Node))
                        return
                      navigate(
                        `/dashboard/agents/${agent.collectionAddress}/${agent.collectionTokenId}`
                      )
                    }}
                  >
                    <TableCell>
                      <AgentAvatar agent={agent} className="size-8" />
                    </TableCell>
                    <TableCell className="font-medium">
                      <div className="flex items-center gap-2">
                        {agent.name}
                        {isProposal && (
                          <Badge variant="outline">Proposal</Badge>
                        )}
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {agent.ensDomain || "—"}
                    </TableCell>
                    <TableCell>
                      <Badge variant="secondary">
                        {agent.attributes["role"]}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {truncateEthAddress(agent.owner)}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {collectionName(agent.collectionAddress)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        {!isProposal &&
                          canManageAgent(agent, address) &&
                          agent.ensDomain && (
                            <Button
                              variant="ghost"
                              size="icon"
                              className="size-6"
                              title="Open agent node"
                              aria-label="Open agent node"
                              onClick={(event) => {
                                event.stopPropagation()
                                openAgentNode(
                                  agent.collectionAddress,
                                  String(agent.collectionTokenId)
                                )
                              }}
                            >
                              <ExternalLink className="size-3.5" />
                            </Button>
                          )}
                        {!isProposal && canManageAgent(agent, address) && (
                          <TransferAgentDialog
                            collectionAddress={agent.collectionAddress}
                            collectionTokenId={String(agent.collectionTokenId)}
                            agentName={agent.name}
                            onTransferred={() => agentsSource.reload()}
                          />
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        )}

        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">
            Page {page} of {totalPages} · {totalElements} {activeFilter.noun}
          </span>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={loading || page <= 1}
              onClick={() => setPage((p) => p - 1)}
            >
              <ChevronLeft />
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={loading || page >= totalPages}
              onClick={() => setPage((p) => p + 1)}
            >
              Next
              <ChevronRight />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
