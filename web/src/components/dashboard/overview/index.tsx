import React from "react"
import { useNavigate } from "react-router-dom"
import { Bot, Hourglass, Layers, MessagesSquare } from "lucide-react"

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
import { useAsyncData } from "@/hooks/use-async-data"
import { formatRelativeTime } from "@/lib/format"
import { listAgents } from "@/queries/api/agents.queries"
import { listAgentCollections } from "@/queries/api/collections.queries"
import { listConversations } from "@/queries/api/conversations.queries"

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

  const mintedAgents = useAsyncData(
    React.useCallback(
      () => listAgents({ page: 1, elementsPerPage: 1, isMinted: true }),
      []
    ),
    "Failed to load agents"
  )
  const proposals = useAsyncData(
    React.useCallback(
      () => listAgents({ page: 1, elementsPerPage: 1, isMinted: false }),
      []
    ),
    "Failed to load proposals"
  )
  const collections = useAsyncData(
    React.useCallback(() => listAgentCollections(), []),
    "Failed to load collections"
  )
  const conversations = useAsyncData(
    React.useCallback(() => listConversations(1, 10), []),
    "Failed to load conversations"
  )

  const conversationList = conversations.data?.data.conversations ?? []
  const recentConversations = conversationList.slice(0, 5)
  const activeConversationsCount = conversationList.filter(
    (conversation) => conversation.status === "active"
  ).length
  const conversationsCount = conversations.data?.data.totalElements ?? 0

  return (
    <div className="flex flex-col gap-6">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Agents"
          value={mintedAgents.data?.data.totalElements ?? 0}
          description="Minted on-chain"
          icon={Bot}
          loading={mintedAgents.loading}
        />
        <StatCard
          title="Proposals"
          value={proposals.data?.data.totalElements ?? 0}
          description="Awaiting mint"
          icon={Hourglass}
          loading={proposals.loading}
        />
        <StatCard
          title="Collections"
          value={collections.data?.data.length ?? 0}
          description="On-chain agent collections"
          icon={Layers}
          loading={collections.loading}
        />
        <StatCard
          title="Active conversations"
          value={activeConversationsCount}
          description={`Of ${conversationsCount} conversations total`}
          icon={MessagesSquare}
          loading={conversations.loading}
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
          {conversations.loading ? (
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
