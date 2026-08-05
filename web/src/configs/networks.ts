import {
  hardhat,
  polygon,
  polygonAmoy,
  type AppKitNetwork,
} from "@reown/appkit/networks"

import { findEvmNetwork } from "@/models/networks"

export const networkByChainId: Record<number, AppKitNetwork> = {
  137: polygon,
  80002: polygonAmoy,
  31337: hardhat,
}

// ENS lives on Ethereum L1: test chains resolve via the Sepolia ENS app,
// production chains via the mainnet one. Chains the console doesn't know (e.g.
// local) have no public ENS app. chainId is the plain id the identity lives on.
export const getEnsAppUrl = (
  chainId: string,
  ensDomain: string
): string | null => {
  const network = findEvmNetwork(`eip155:${chainId}`)
  if (!network) return null
  const base = network.testnet
    ? "https://sepolia.app.ens.domains"
    : "https://app.ens.domains"
  return `${base}/${ensDomain}`
}

// Block-explorer address URL for a chain, or null if the chain has no public
// explorer (e.g. local hardhat). The base comes from the AppKit network config
// so it stays in sync with the chains the app actually supports.
export const getScannerAddressUrl = (
  chainId: string,
  address: string
): string | null => {
  const base = networkByChainId[Number(chainId)]?.blockExplorers?.default?.url
  if (!base || !address) return null
  return `${base}/address/${address}`
}

// Marketplace path slugs per chain id. Only chains actually listed on the
// marketplace appear here; others return null so the link is simply hidden.
// Testnets resolve via a separate marketplace host (see helpers below).
const MARKETPLACE_SLUGS: Record<string, { openSea: string; rarible: string }> =
  {
    "1": { openSea: "ethereum", rarible: "ethereum" },
    "137": { openSea: "matic", rarible: "polygon" },
    "8453": { openSea: "base", rarible: "base" },
    "42161": { openSea: "arbitrum", rarible: "arbitrum" },
    "80002": { openSea: "amoy", rarible: "polygon" }, // Polygon Amoy (testnet)
  }

// OpenSea item URL for a minted agent, or null if the chain isn't on OpenSea or
// the agent isn't minted. chainId is the plain id the identity lives on.
export const getOpenSeaUrl = (
  chainId: string,
  collectionAddress: string,
  tokenId: number | null
): string | null => {
  const slugs = MARKETPLACE_SLUGS[chainId]
  const network = findEvmNetwork(`eip155:${chainId}`)
  if (!slugs || !network || tokenId == null || !collectionAddress) return null
  const base = network.testnet
    ? "https://testnets.opensea.io"
    : "https://opensea.io"
  return `${base}/item/${slugs.openSea}/${collectionAddress}/${tokenId}`
}

// Rarible token URL for a minted agent, or null if the chain isn't on Rarible or
// the agent isn't minted. chainId is the plain id the identity lives on.
export const getRaribleUrl = (
  chainId: string,
  collectionAddress: string,
  tokenId: number | null
): string | null => {
  const slugs = MARKETPLACE_SLUGS[chainId]
  const network = findEvmNetwork(`eip155:${chainId}`)
  if (!slugs || !network || tokenId == null || !collectionAddress) return null
  const base = network.testnet
    ? "https://testnet.rarible.com"
    : "https://rarible.com"
  return `${base}/token/${slugs.rarible}/${collectionAddress}:${tokenId}`
}
