import { config } from "@/configs/env"

// Turn a stored image value into a displayable URL. Values are usually bare IPFS
// CIDs (e.g. the agent's "default" image); if a value is already a full URI
// (ipfs://…, http(s)://… or data:) it's passed through unchanged.
export const ipfsImageUrl = (value?: string | null): string | null => {
  if (!value) return null
  if (/^(https?:|ipfs:|data:)/.test(value)) return value
  return `${config.ipfsGateway}${value}`
}
