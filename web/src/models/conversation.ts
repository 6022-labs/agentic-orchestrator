export type ConversationStatus = "active" | "concluded"

export type Conversation = {
  id: string
  bridgeId: string
  description?: string
  status: ConversationStatus
  concludedAt?: string
  createdAt: string
  updatedAt: string
}

export type ListConversationsResponse = {
  page: number
  elementsPerPage: number
  totalElements: number
  conversations: Conversation[]
}

export type ConversationMessage = {
  id: string
  conversationId: string
  participantId: string
  content: string
  createdAt: string
}

export type ConversationParticipant = {
  id: string
  collectionAddress?: string
  collectionTokenId?: number
  ensDomain?: string
  role?: string
  displayName?: string
}
