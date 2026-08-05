import { BrowserProvider, ethers, type Eip1193Provider } from "ethers"

import { AgentCollectionCreatorABI } from "@/abi/agent-collection-creator"
import { getPrimaryEvmChainConfig } from "@/configs/app-config"
import { getScaledTxGasOverrides } from "@/lib/tx-gas-overrides"

const collectionCreatedInterface = new ethers.Interface([
  "event CollectionCreated(uint256 indexed collectionId, address collectionAddress, string name, string symbol)",
])

const logAndRethrow = (error: unknown) => {
  try {
    const iface = new ethers.Interface(AgentCollectionCreatorABI)
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const parsed = iface.parseError((error as any)?.data)
    if (parsed?.name) {
      ;(error as { parsedRevertName?: string }).parsedRevertName = parsed.name
    }
    console.error("AgentCollectionCreator call failed", error, parsed)
  } catch {
    console.error("AgentCollectionCreator call failed", error)
  }
  throw error
}

// Creates the collection (clone + init + manager registration in one tx) and
// returns the new collection address parsed from the manager's
// CollectionCreated log, or null if the log can't be decoded.
export const createCollection = async (
  name: string,
  symbol: string,
  description: string,
  admin: string,
  moderator: string,
  walletProvider: Eip1193Provider
): Promise<string | null> => {
  const chainConfig = getPrimaryEvmChainConfig()
  const provider = new BrowserProvider(walletProvider)
  const signer = await provider.getSigner()
  const contract = new ethers.Contract(
    chainConfig.collectionCreatorAddress,
    AgentCollectionCreatorABI,
    signer
  )

  try {
    const tx = await contract.createCollection(
      name,
      symbol,
      description,
      admin,
      moderator,
      await getScaledTxGasOverrides(signer)
    )
    const receipt = await tx.wait()
    const managerAddress = chainConfig.collectionsManagerAddress.toLowerCase()
    for (const log of receipt?.logs ?? []) {
      if (log.address.toLowerCase() !== managerAddress) continue
      try {
        const parsed = collectionCreatedInterface.parseLog(log)
        if (parsed?.name === "CollectionCreated") {
          return parsed.args.collectionAddress as string
        }
      } catch {
        // Not a CollectionCreated log; keep scanning.
      }
    }
    return null
  } catch (error) {
    return logAndRethrow(error)
  }
}
