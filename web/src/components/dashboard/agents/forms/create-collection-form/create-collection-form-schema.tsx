import { z } from "zod"

import { ensLabelSchema } from "@/lib/ens-label"

const addressRegex = /^0x[0-9a-fA-F]{40}$/

export const createCollectionFormSchema = z.object({
  // The contract requires a valid ENS label.
  name: ensLabelSchema,
  symbol: z.string().min(1, "Symbol is required"),
  description: z.string().optional(),
  admin: z.string().regex(addressRegex, "Must be a 0x-prefixed address"),
  moderator: z.string().regex(addressRegex, "Must be a 0x-prefixed address"),
})

export type CreateCollectionFormSchemaType = z.infer<
  typeof createCollectionFormSchema
>
