import { z } from "zod"

export const ENS_LABEL_MAX_LENGTH = 63

export function normalizeEnsLabel(value: string): string {
  return value.toLowerCase().replace(/\s+/g, "-")
}

export function isValidEnsLabel(value: string): boolean {
  if (value.length === 0 || value.length > ENS_LABEL_MAX_LENGTH) return false
  if (value.startsWith("-") || value.endsWith("-")) return false
  if (value.includes("--")) return false
  return /^[a-z0-9-]+$/.test(value)
}

export const ensLabelSchema = z
  .string()
  .refine(
    isValidEnsLabel,
    "Lowercase letters, digits and hyphens; no leading, trailing or repeated hyphens; max 63 chars"
  )
