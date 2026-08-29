import { getAppConfig } from "@/configs/app-config"

// Used when the backend serves no list, so images still resolve.
const FALLBACK_GATEWAY_URL_TEMPLATES = [
  "https://w3s.link/ipfs/{cid}",
  "https://4everland.io/ipfs/{cid}",
  "https://dweb.link/ipfs/{cid}",
  "https://ipfs.io/ipfs/{cid}",
]

// Candidate URLs to try in order; an already-displayable URI is the only candidate.
export const ipfsImageUrls = (value?: string | null): string[] => {
  if (!value) return []

  const trimmed = value.trim()
  if (!trimmed) return []
  if (/^(https?:|data:)/.test(trimmed)) return [trimmed]

  // A path suffix can never go into a subdomain template's host.
  const [cid, ...segments] = trimmed.replace(/^ipfs:\/\//, "").split("/")
  const path = segments.length > 0 ? `/${segments.join("/")}` : ""

  const configured = getAppConfig().ipfsGatewayUrlTemplates
  const templates = configured?.length
    ? configured
    : FALLBACK_GATEWAY_URL_TEMPLATES

  return templates.map((template) =>
    template.includes("{cid}")
      ? template.replaceAll("{cid}", cid) + path
      : template + cid + path
  )
}
