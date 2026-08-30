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

export type PromptTemplateKey = (typeof PROMPT_TEMPLATE_KEYS)[number]

export const CLONE_PROMPT_TEMPLATE_KEYS: PromptTemplateKey[] = [
  "clone-system",
  "clone-default",
]

export const EXPERT_PROMPT_TEMPLATE_KEYS: PromptTemplateKey[] = [
  "expert-system",
  "expert-default",
]

export const FACILITATOR_PROMPT_TEMPLATE_KEYS: PromptTemplateKey[] = [
  "facilitator-system",
  "facilitator-default",
  "facilitator-conclude",
]

// Free-form/unknown roles (incl. human and empty) fall back to expert.
export const getPromptTemplateKeysForRole = (
  role: string
): PromptTemplateKey[] => {
  switch (role.trim().toLowerCase()) {
    case "facilitator":
      return FACILITATOR_PROMPT_TEMPLATE_KEYS
    case "clone":
      return CLONE_PROMPT_TEMPLATE_KEYS
    default:
      return EXPERT_PROMPT_TEMPLATE_KEYS
  }
}
