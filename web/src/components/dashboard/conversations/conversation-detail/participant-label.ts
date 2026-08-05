import type { ConversationParticipant } from "@/models/conversation"

export const participantLabel = (participant?: ConversationParticipant) =>
  participant?.displayName ||
  participant?.ensDomain ||
  participant?.id ||
  "Unknown participant"
