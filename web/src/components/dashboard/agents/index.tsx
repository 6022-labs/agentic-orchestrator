import React from "react"
import { toast } from "sonner"
import { useNavigate } from "react-router-dom"
import { Bot, ChevronLeft, ChevronRight } from "lucide-react"

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
import AgentAvatar from "@/components/dashboard/agents/agent-avatar"
import { extractErrorMessage } from "@/lib/api-errors"
import { truncateEthAddress } from "@/lib/utils"
import type { Agent } from "@/models/agent"
import type { AgentCollection } from "@/models/agent-collection"
import { listAgents } from "@/queries/api/agents.queries"
import { listAgentCollections } from "@/queries/api/collections.queries"

const ELEMENTS_PER_PAGE = 20

export default function AgentsSection() {
  const navigate = useNavigate()
  const [loading, setLoading] = React.useState(true)
  const [agents, setAgents] = React.useState<Agent[]>([])
  const [collections, setCollections] = React.useState<AgentCollection[]>([])
  const [page, setPage] = React.useState(1)
  const [totalElements, setTotalElements] = React.useState(0)

  const load = React.useCallback(async () => {
    setLoading(true)
    try {
      const [agentsResponse, collectionsResponse] = await Promise.all([
        listAgents({ page, elementsPerPage: ELEMENTS_PER_PAGE }),
        listAgentCollections(),
      ])
      setAgents(agentsResponse.data.agents ?? [])
      setTotalElements(agentsResponse.data.totalElements)
      setCollections(collectionsResponse.data ?? [])
    } catch (error) {
      toast.error(extractErrorMessage(error, "Failed to load agents"))
    } finally {
      setLoading(false)
    }
  }, [page])

  React.useEffect(() => {
    void load()
  }, [load])

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
          <CreateCollectionDialog onCreated={() => void load()} />
        </div>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
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
            description="Registered agents show up here once they are minted."
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
              </TableRow>
            </TableHeader>
            <TableBody>
              {agents.map((agent) => (
                <TableRow
                  key={agent.id}
                  className="cursor-pointer"
                  onClick={() =>
                    navigate(
                      `/dashboard/agents/${agent.collectionAddress}/${agent.collectionTokenId}`
                    )
                  }
                >
                  <TableCell>
                    <AgentAvatar agent={agent} className="size-8" />
                  </TableCell>
                  <TableCell className="font-medium">{agent.name}</TableCell>
                  <TableCell className="text-muted-foreground">
                    {agent.ensDomain || "—"}
                  </TableCell>
                  <TableCell>
                    <Badge variant="secondary">{agent.attributes["role"]}</Badge>
                  </TableCell>
                  <TableCell className="text-muted-foreground">
                    {truncateEthAddress(agent.owner)}
                  </TableCell>
                  <TableCell className="text-muted-foreground">
                    {collectionName(agent.collectionAddress)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}

        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">
            Page {page} of {totalPages} · {totalElements} agents
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
