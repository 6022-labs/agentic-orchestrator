import React from "react"
import { toast } from "sonner"
import { useNavigate } from "react-router-dom"
import { Bot, Layers, MessagesSquare } from "lucide-react"

import { Badge } from "@/components/ui/badge"
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
import { extractErrorMessage } from "@/lib/api-errors"
import { formatRelativeTime } from "@/lib/format"
import type { Conversation } from "@/models/conversation"
import { listAgents } from "@/queries/api/agents.queries"
import { listAgentCollections } from "@/queries/api/collections.queries"
import { listConversations } from "@/queries/api/conversations.queries"

type OverviewStats = {
  agentsCount: number
  collectionsCount: number
  conversationsCount: number
  activeConversationsCount: number
}

function StatCard({
  title,
  value,
  description,
  icon: Icon,
  loading,
}: {
  title: string
  value: number
  description: string
  icon: typeof Bot
  loading: boolean
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className="size-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        {loading ? (
          <Skeleton className="h-8 w-16" />
        ) : (
          <div className="text-2xl font-bold">{value}</div>
        )}
        <p className="text-xs text-muted-foreground">{description}</p>
      </CardContent>
    </Card>
  )
}

export default function OverviewSection() {
  const navigate = useNavigate()
  const [loading, setLoading] = React.useState(true)
  const [stats, setStats] = React.useState<OverviewStats | null>(null)
  const [recentConversations, setRecentConversations] = React.useState<
    Conversation[]
  >([])

  const load = React.useCallback(async () => {
    setLoading(true)
    try {
      const [agentsResponse, collectionsResponse, conversationsResponse] =
        await Promise.all([
          listAgents({ page: 1, elementsPerPage: 1 }),
          listAgentCollections(),
          listConversations(1, 10),
        ])

      const conversations = conversationsResponse.data.conversations ?? []
      setStats({
        agentsCount: agentsResponse.data.totalElements,
        collectionsCount: collectionsResponse.data.length,
        conversationsCount: conversationsResponse.data.totalElements,
        activeConversationsCount: conversations.filter(
          (c) => c.status === "active"
        ).length,
      })
      setRecentConversations(conversations.slice(0, 5))
    } catch (error) {
      toast.error(extractErrorMessage(error, "Failed to load overview"))
    } finally {
      setLoading(false)
    }
  }, [])

  React.useEffect(() => {
    void load()
  }, [load])

  return (
    <div className="flex flex-col gap-6">
      <div className="grid gap-4 md:grid-cols-3">
        <StatCard
          title="Agents"
          value={stats?.agentsCount ?? 0}
          description="Registered in the orchestrator"
          icon={Bot}
          loading={loading}
        />
        <StatCard
          title="Collections"
          value={stats?.collectionsCount ?? 0}
          description="On-chain agent collections"
          icon={Layers}
          loading={loading}
        />
        <StatCard
          title="Active conversations"
          value={stats?.activeConversationsCount ?? 0}
          description={`Of ${stats?.conversationsCount ?? 0} conversations total`}
          icon={MessagesSquare}
          loading={loading}
        />
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Recent conversations</CardTitle>
          <CardDescription>
            The latest conversations the orchestrator bridged.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="space-y-2">
              <Skeleton className="h-9 w-full" />
              <Skeleton className="h-9 w-full" />
              <Skeleton className="h-9 w-full" />
            </div>
          ) : recentConversations.length === 0 ? (
            <EmptyState
              icon={MessagesSquare}
              title="No conversations yet"
              description="Conversations appear here once agents start talking."
            />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Description</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Updated</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {recentConversations.map((conversation) => (
                  <TableRow
                    key={conversation.id}
                    className="cursor-pointer"
                    onClick={() =>
                      navigate(`/dashboard/conversations/${conversation.id}`)
                    }
                  >
                    <TableCell className="max-w-96 truncate">
                      {conversation.description || conversation.bridgeId}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          conversation.status === "active"
                            ? "default"
                            : "secondary"
                        }
                      >
                        {conversation.status}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {formatRelativeTime(conversation.updatedAt)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
