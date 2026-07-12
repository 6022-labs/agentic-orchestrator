import axios from "axios"

// EvmChainConfig is one EVM chain's config surfaced by the gateway.
export type EvmChainConfig = {
  collectionsManagerAddress: string
  collectionCreatorAddress: string
}

// Config served by the gateway's GET /config at boot; evmChains keyed by chain id.
export type AppConfig = {
  evmChains: Record<string, EvmChainConfig>
}

let appConfig: AppConfig | null = null
let pending: Promise<AppConfig> | null = null

// Idempotent: concurrent/repeated calls (StrictMode double-mount) share one request.
export async function loadAppConfig(): Promise<AppConfig> {
  if (appConfig) return appConfig
  if (!pending) {
    pending = axios
      .get<AppConfig>("/config")
      .then(({ data }) => {
        appConfig = data
        return data
      })
      .finally(() => {
        pending = null
      })
  }
  return pending
}

export function getAppConfig(): AppConfig {
  if (!appConfig) {
    throw new Error("App config read before loadAppConfig() resolved")
  }
  return appConfig
}

// getEvmChainConfig returns the config for a chain id, or undefined if the
// gateway does not serve that chain.
function getEvmChainConfig(
  chainId: string,
): EvmChainConfig | undefined {
  return getAppConfig().evmChains[chainId]
}

const env = import.meta.env

// getPrimaryEvmChainConfig resolves the single chain the console operates on,
// falling back to the Vite env so a misconfigured gateway doesn't brick it.
export function getPrimaryEvmChainConfig(): EvmChainConfig {
  const chainIds = Object.keys(getAppConfig().evmChains ?? {})
  const configured =
    chainIds.length > 0 ? getEvmChainConfig(chainIds[0]) : undefined
  return {
    collectionsManagerAddress:
      configured?.collectionsManagerAddress ||
      env.VITE_AGENT_COLLECTIONS_MANAGER_ADDRESS ||
      "",
    collectionCreatorAddress:
      configured?.collectionCreatorAddress ||
      env.VITE_AGENT_COLLECTION_CREATOR_ADDRESS ||
      "",
  }
}
