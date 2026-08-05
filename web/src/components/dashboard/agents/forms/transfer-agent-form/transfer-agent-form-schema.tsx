import { z } from "zod"

const addressRegex = /^0x[0-9a-fA-F]{40}$/

export const transferAgentFormSchema = z.object({
  to: z.string().regex(addressRegex, "Must be a 0x-prefixed address"),
})

export type TransferAgentFormSchemaType = z.infer<
  typeof transferAgentFormSchema
>
