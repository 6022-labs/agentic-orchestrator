import React from "react"
import { toast } from "sonner"
import { useNavigate } from "react-router-dom"
import { ChevronLeft, ChevronRight, MessagesSquare } from "lucide-react"

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
import { extractErrorMessage } from "@/lib/api-errors"
import { formatRelativeTime } from "@/lib/format"
import type { Conversation } from "@/models/conversation"
import { listConversations } from "@/queries/api/conversations.queries"

const ELEMENTS_PER_PAGE = 20

export default function ConversationsSection() {
  const navigate = useNavigate()
  const [loading, setLoading] = React.useState(true)
  const [conversations, setConversations] = React.useState<Conversation[]>([])
  const [page, setPage] = React.useState(1)
  const [totalElements, setTotalElements] = React.useState(0)

  const load = React.useCallback(async () => {
    setLoading(true)
    try {
      const { data } = await listConversations(page, ELEMENTS_PER_PAGE)
      setConversations(data.conversations ?? [])
      setTotalElements(data.totalElements)
    } catch (error) {
      toast.error(extractErrorMessage(error, "Failed to load conversations"))
    } finally {
      setLoading(false)
    }
  }, [page])

  React.useEffect(() => {
    void load()
  }, [load])

  const totalPages = Math.max(Math.ceil(totalElements / ELEMENTS_PER_PAGE), 1)

  return (
    <Card>
      <CardHeader>
        <CardTitle>Conversations</CardTitle>
        <CardDescription>
          Every conversation the orchestrator bridged between agents.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        {loading ? (
          <div className="space-y-2">
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
          </div>
        ) : conversations.length === 0 ? (
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
              {conversations.map((conversation) => (
                <TableRow
                  key={conversation.id}
                  className="cursor-pointer"
                  onClick={() =>
                    navigate(`/dashboard/conversations/${conversation.id}`)
                  }
                >
                  <TableCell className="max-w-96 truncate font-medium">
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
                  <TableCell className="text-muted-foreground">
                    {formatRelativeTime(conversation.updatedAt)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}

        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">
            Page {page} of {totalPages} · {totalElements} conversations
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
