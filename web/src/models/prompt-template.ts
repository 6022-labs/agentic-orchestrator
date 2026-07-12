export type PromptTemplate = {
  collectionAddress: string
  collectionTokenId: number
  key: string
  content: string
}

export const PROMPT_TEMPLATE_KEYS = [
  "clone-system",
  "clone-default",
  "expert-system",
  "expert-default",
  "facilitator-system",
  "facilitator-default",
  "facilitator-conclude",
] as const

