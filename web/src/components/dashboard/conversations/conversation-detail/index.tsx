import React from "react"
import { toast } from "sonner"
import { Link, useParams } from "react-router-dom"
import { ArrowLeft, Gavel, MessagesSquare, Users } from "lucide-react"

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
import ConfirmDialog from "@/components/dashboard/confirm-dialog"
import EmptyState from "@/components/dashboard/empty-state"
import { useAsyncData } from "@/hooks/use-async-data"
import { extractErrorMessage } from "@/lib/api-errors"
import { formatRelativeTime } from "@/lib/format"
import {
  forceConcludeConversation,
  getConversation,
  getConversationMessages,
  getConversationParticipants,
} from "@/queries/api/conversations.queries"
import ConversationTranscript from "@/components/dashboard/conversations/conversation-detail/conversation-transcript"
import { participantLabel } from "@/components/dashboard/conversations/conversation-detail/participant-label"

export default function ConversationDetailSection() {
  const { conversationId } = useParams<{ conversationId: string }>()
  const id = conversationId ?? ""

  const [concluding, setConcluding] = React.useState(false)

  const conversationSource = useAsyncData(
    React.useCallback(() => getConversation(id), [id]),
    "Failed to load the conversation"
  )
  const messagesSource = useAsyncData(
    React.useCallback(() => getConversationMessages(id), [id]),
    "Failed to load the transcript"
  )
  const participantsSource = useAsyncData(
    React.useCallback(() => getConversationParticipants(id), [id]),
    "Failed to load the participants"
  )

  const conversation = conversationSource.data?.data ?? null
  const messages = messagesSource.data?.data ?? []
  const participants = participantsSource.data?.data ?? []

  const onForceConclude = async () => {
    setConcluding(true)
    try {
      await forceConcludeConversation(id)
      toast.success("Conversation concluded")
      conversationSource.reload()
      messagesSource.reload()
      participantsSource.reload()
    } catch (error) {
      toast.error(
        extractErrorMessage(error, "Failed to conclude the conversation")
      )
    } finally {
      setConcluding(false)
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <Button variant="ghost" size="sm" asChild>
          <Link to="/dashboard/conversations">
            <ArrowLeft />
            Back to conversations
          </Link>
        </Button>
        {conversation?.status === "active" ? (
          <ConfirmDialog
            title="Force conclude this conversation?"
            description="The conversation is closed for every participant. This cannot be undone."
            confirmLabel="Force conclude"
            onConfirm={() => void onForceConclude()}
          >
            <Button variant="destructive" size="sm" loading={concluding}>
              <Gavel />
              Force conclude
            </Button>
          </ConfirmDialog>
        ) : null}
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {conversationSource.loading ? (
          <Skeleton className="h-56 w-full" />
        ) : !conversation ? (
          <EmptyState
            icon={MessagesSquare}
            title="Conversation not found"
            description="The orchestrator has no conversation with this id."
          />
        ) : (
          <Card>
            <CardHeader>
              <CardTitle>Conversation</CardTitle>
              <CardDescription>
                {conversation.description || conversation.bridgeId}
              </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-col gap-3 text-sm">
              <div className="flex items-center gap-2">
                <span className="text-muted-foreground">Status</span>
                <Badge
                  variant={
                    conversation.status === "active" ? "default" : "secondary"
                  }
                >
                  {conversation.status}
                </Badge>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-muted-foreground">Created</span>
                <span>{formatRelativeTime(conversation.createdAt)}</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-muted-foreground">Updated</span>
                <span>{formatRelativeTime(conversation.updatedAt)}</span>
              </div>
              {conversation.concludedAt ? (
                <div className="flex items-center gap-2">
                  <span className="text-muted-foreground">Concluded</span>
                  <span>{formatRelativeTime(conversation.concludedAt)}</span>
                </div>
              ) : null}
            </CardContent>
          </Card>
        )}

        <Card>
          <CardHeader>
            <CardTitle>Participants</CardTitle>
            <CardDescription>
              Everyone taking part in the conversation.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {participantsSource.loading ? (
              <div className="space-y-2">
                <Skeleton className="h-5 w-full" />
                <Skeleton className="h-5 w-full" />
              </div>
            ) : participants.length === 0 ? (
              <EmptyState icon={Users} title="No participants" />
            ) : (
              <ul className="flex flex-col gap-2">
                {participants.map((participant) => (
                  <li
                    key={participant.id}
                    className="flex items-center gap-2 text-sm"
                  >
                    <span className="truncate font-medium">
                      {participantLabel(participant)}
                    </span>
                    {participant.role ? (
                      <Badge variant="outline">{participant.role}</Badge>
                    ) : null}
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
      </div>

      {messagesSource.loading ? (
        <Skeleton className="h-64 w-full" />
      ) : (
        <ConversationTranscript
          messages={messages}
          participants={participants}
        />
      )}
    </div>
  )
}
