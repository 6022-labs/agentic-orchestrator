import { BrowserProvider, type Eip1193Provider, ethers } from "ethers"

import { AgentCollectionsManagerABI } from "@/abi/agent-collections-manager"
import { getPrimaryEvmChainConfig } from "@/configs/app-config"

// moderatorOf lists the collection addresses the connected wallet moderates.
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

  return await contract.moderatorOf(await connectedSigner.getAddress())
}
