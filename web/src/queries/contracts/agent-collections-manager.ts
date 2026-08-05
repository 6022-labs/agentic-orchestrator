import { BrowserProvider, type Eip1193Provider, ethers } from "ethers"

import { AgentCollectionsManagerABI } from "@/abi/agent-collections-manager"
import { getPrimaryEvmChainConfig } from "@/configs/app-config"

// moderatorOf lists the collection addresses the connected wallet moderates,
// walking the paginated on-chain cursor until exhausted.
export async function moderatorOf(
  walletProvider: Eip1193Provider
): Promise<string[]> {
  const provider = new BrowserProvider(walletProvider)
  const connectedSigner = await provider.getSigner()

  const contract = new ethers.Contract(
    getPrimaryEvmChainConfig().collectionsManagerAddress,
    AgentCollectionsManagerABI,
    connectedSigner
  )

  const wallet = await connectedSigner.getAddress()
  const pageSize = 50n
  const collections: string[] = []
  let offset = 0n
  do {
    const [page, next] = await contract.moderatorOf(wallet, offset, pageSize)
    collections.push(...page)
    offset = next
  } while (offset !== 0n)

  return collections
}
