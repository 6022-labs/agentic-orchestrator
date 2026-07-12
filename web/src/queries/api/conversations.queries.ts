import axios from "axios"

import type {
  Conversation,
  ConversationMessage,
  ConversationParticipant,
  ListConversationsResponse,
} from "@/models/conversation"

// GET /conversations — paginated conversations, most recent first.
export const listConversations = async (
  page?: number,
  elementsPerPage?: number
) => {
  return await axios.get<ListConversationsResponse>("/conversations", {
    params: { page, elementsPerPage },
  })
}

// GET /conversations/{id} — one conversation.
export const getConversation = async (conversationId: string) => {
  return await axios.get<Conversation>(`/conversations/${conversationId}`)
}

// GET /conversations/{id}/messages — one conversation's transcript.
export const getConversationMessages = async (conversationId: string) => {
  return await axios.get<ConversationMessage[]>(
    `/conversations/${conversationId}/messages`
  )
}

// GET /conversations/{id}/participants — everyone in the conversation.
export const getConversationParticipants = async (conversationId: string) => {
  return await axios.get<ConversationParticipant[]>(
    `/conversations/${conversationId}/participants`
  )
}

// POST /conversations/{id}/force-conclude — moderator override, 204 on success.
export const forceConcludeConversation = async (conversationId: string) => {
  return await axios.post<void>(
    `/conversations/${conversationId}/force-conclude`
  )
}
