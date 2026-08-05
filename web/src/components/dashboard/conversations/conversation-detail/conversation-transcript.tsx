import React from "react"
import { MessagesSquare } from "lucide-react"

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import EmptyState from "@/components/dashboard/empty-state"
import { participantLabel } from "@/components/dashboard/conversations/conversation-detail/participant-label"
import { formatRelativeTime } from "@/lib/format"
import type {
  ConversationMessage,
  ConversationParticipant,
} from "@/models/conversation"

export default function ConversationTranscript({
  messages,
  participants,
}: {
  messages: ConversationMessage[]
  participants: ConversationParticipant[]
}) {
  const participantById = React.useMemo(() => {
    const map = new Map<string, ConversationParticipant>()
    for (const participant of participants) {
      map.set(participant.id, participant)
    }
    return map
  }, [participants])

  return (
    <Card>
      <CardHeader>
        <CardTitle>Transcript</CardTitle>
        <CardDescription>
          {messages.length} message{messages.length === 1 ? "" : "s"}
        </CardDescription>
      </CardHeader>
      <CardContent>
        {messages.length === 0 ? (
          <EmptyState
            icon={MessagesSquare}
            title="No messages yet"
            description="The transcript is empty."
          />
        ) : (
          <div className="space-y-3 rounded-xl border bg-muted/30 p-4">
            {messages.map((message, index) => {
              const participant = participantById.get(message.participantId)
              const previous = messages[index - 1]
              const sameAuthor =
                previous?.participantId === message.participantId

              return (
                <div key={message.id} className="flex flex-col gap-1">
                  {!sameAuthor ? (
                    <span className="text-xs font-medium text-muted-foreground">
                      {participantLabel(participant)}
                    </span>
                  ) : null}
                  <div className="max-w-[85%] rounded-2xl border bg-background px-3 py-2 text-sm whitespace-pre-wrap">
                    {message.content}
                  </div>
                  <span className="text-xs text-muted-foreground">
                    {formatRelativeTime(message.createdAt)}
                  </span>
                </div>
              )
            })}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
