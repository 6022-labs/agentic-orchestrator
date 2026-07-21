import { BrowserProvider, type Eip1193Provider, ethers } from "ethers"

import { AgentCollectionABI } from "@/abi/agent-collection"
import type { MintProposal } from "@/models/mint-proposal"
import type { AddOrUpdateImageProposal } from "@/models/add-or-update-image-proposal"

const collectionContract = async (
  collectionAddress: string,
  walletProvider: Eip1193Provider
) => {
  const provider = new BrowserProvider(walletProvider)
  const connectedSigner = await provider.getSigner()
  return new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    connectedSigner
  )
}

const logDecodedError = (error: unknown) => {
  const iface = new ethers.Interface(AgentCollectionABI)
  const data = (error as { data?: string })?.data
  try {
    console.error(error, data ? iface.parseError(data) : undefined)
  } catch {
    console.error(error)
  }
}

// hasModeratorRole checks whether an address holds MODERATOR_ROLE on a collection.
export async function hasModeratorRole(
  collectionAddress: string,
  address: string,
  walletProvider: Eip1193Provider
): Promise<boolean> {
  const contract = await collectionContract(collectionAddress, walletProvider)
  const moderatorRole = await contract.MODERATOR_ROLE()
  return await contract.hasRole(moderatorRole, address)
}

export async function getMintProposals(
  collectionAddress: string,
  walletProvider: Eip1193Provider
): Promise<MintProposal[]> {
  const contract = await collectionContract(collectionAddress, walletProvider)

  try {
    const mintProposalsLength = await contract.mintProposalsLength()
    const promises = []

    for (let i = 0; i < mintProposalsLength; i++) {
      promises.push(contract.mintProposal(i))
    }

    const mintProposals = await Promise.all(promises)

    return mintProposals.map((proposal) => ({
      id: proposal.id,
      collectionAddress: collectionAddress,
      name: proposal.name,
      role:
        proposal.attributes?.find(
          (attribute: { key: string }) => attribute.key === "role"
        )?.value ?? "",
      owner: proposal.owner,
      mission: proposal.mission,
      creator: proposal.creator,
      wallets: proposal.wallets,
      images: proposal.images,
    }))
  } catch (error) {
    logDecodedError(error)
    throw error
  }
}

export async function mintFromProposal(
  collectionAddress: string,
  proposalIndex: number,
  walletProvider: Eip1193Provider
) {
  const contract = await collectionContract(collectionAddress, walletProvider)

  try {
    const tx = await contract.mintFromProposal(proposalIndex)
    await tx.wait()
    return tx
  } catch (error) {
    logDecodedError(error)
    throw error
  }
}

export async function refuseMintProposal(
  collectionAddress: string,
  proposalIndex: number,
  walletProvider: Eip1193Provider
) {
  const contract = await collectionContract(collectionAddress, walletProvider)

  try {
    const tx = await contract.refuseMintProposal(proposalIndex)
    await tx.wait()
  } catch (error) {
    logDecodedError(error)
    throw error
  }
}

export async function getAddOrUpdateAgentImageProposals(
  collectionAddress: string,
  walletProvider: Eip1193Provider
): Promise<AddOrUpdateImageProposal[]> {
  const contract = await collectionContract(collectionAddress, walletProvider)

  try {
    const proposalsLength =
      await contract.addOrUpdateAgentImageProposalsLength()
    const promises = []

    for (let i = 0; i < proposalsLength; i++) {
      promises.push(contract.addOrUpdateAgentImageProposal(i))
    }

    const proposals = await Promise.all(promises)

    return proposals.map((proposal) => ({
      id: proposal.id,
      collectionAddress: collectionAddress,
      tokenId: proposal.tokenId,
      image: {
        key: proposal.image.key,
        value: proposal.image.value,
      },
    }))
  } catch (error) {
    logDecodedError(error)
    throw error
  }
}

export async function addOrUpdateAgentImageFromProposal(
  collectionAddress: string,
  proposalIndex: number,
  walletProvider: Eip1193Provider
) {
  const contract = await collectionContract(collectionAddress, walletProvider)

  try {
    const tx = await contract.addOrUpdateAgentImageFromProposal(proposalIndex)
    await tx.wait()
  } catch (error) {
    logDecodedError(error)
    throw error
  }
}

export async function refuseAddOrUpdateAgentImageProposal(
  collectionAddress: string,
  proposalIndex: number,
  walletProvider: Eip1193Provider
) {
  const contract = await collectionContract(collectionAddress, walletProvider)

  try {
    const tx = await contract.refuseAddOrUpdateAgentImageProposal(proposalIndex)
    await tx.wait()
  } catch (error) {
    logDecodedError(error)
    throw error
  }
}

// transferFrom moves the agent NFT from the connected wallet to another owner.
export async function transferFrom(
  collectionAddress: string,
  tokenId: string,
  to: string,
  walletProvider: Eip1193Provider
): Promise<void> {
  const provider = new BrowserProvider(walletProvider)
  const signer = await provider.getSigner()
  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    signer
  )

  try {
    const tx = await contract.transferFrom(await signer.getAddress(), to, tokenId)
    await tx.wait()
  } catch (error) {
    logDecodedError(error)
    throw error
  }
}
