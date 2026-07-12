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
import { extractErrorMessage } from "@/lib/api-errors"
import { formatRelativeTime } from "@/lib/format"
import type {
  Conversation,
  ConversationMessage,
  ConversationParticipant,
} from "@/models/conversation"
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

  const [loading, setLoading] = React.useState(true)
  const [conversation, setConversation] = React.useState<Conversation | null>(
    null
  )
  const [messages, setMessages] = React.useState<ConversationMessage[]>([])
  const [participants, setParticipants] = React.useState<
    ConversationParticipant[]
  >([])
  const [concluding, setConcluding] = React.useState(false)

  const load = React.useCallback(async () => {
    if (!conversationId) return

    setLoading(true)
    try {
      const [conversationResponse, messagesResponse, participantsResponse] =
        await Promise.all([
          getConversation(conversationId),
          getConversationMessages(conversationId),
          getConversationParticipants(conversationId),
        ])
      setConversation(conversationResponse.data)
      setMessages(messagesResponse.data ?? [])
      setParticipants(participantsResponse.data ?? [])
    } catch (error) {
      toast.error(extractErrorMessage(error, "Failed to load the conversation"))
    } finally {
      setLoading(false)
    }
  }, [conversationId])

  React.useEffect(() => {
    void load()
  }, [load])

  const onForceConclude = async () => {
    if (!conversationId) return

    setConcluding(true)
    try {
      await forceConcludeConversation(conversationId)
      toast.success("Conversation concluded")
      await load()
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

      {loading ? (
        <div className="space-y-4">
          <Skeleton className="h-32 w-full" />
          <Skeleton className="h-64 w-full" />
        </div>
      ) : !conversation ? (
        <EmptyState
          icon={MessagesSquare}
          title="Conversation not found"
          description="The orchestrator has no conversation with this id."
        />
      ) : (
        <>
          <div className="grid gap-6 lg:grid-cols-2">
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

            <Card>
              <CardHeader>
                <CardTitle>Participants</CardTitle>
                <CardDescription>
                  Everyone taking part in the conversation.
                </CardDescription>
              </CardHeader>
              <CardContent>
                {participants.length === 0 ? (
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

          <ConversationTranscript
            messages={messages}
            participants={participants}
          />
        </>
      )}
    </div>
  )
}
