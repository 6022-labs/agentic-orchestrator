import { MintProposal } from "@/models/mint-proposal";
import { AgentCollectionABI } from "@/abi/agent-collection";
import { BrowserProvider, Eip1193Provider, ethers } from "ethers";
import { AddOrUpdateImageProposal } from "@/models/add-or-update-image-proposal";
import { mapToKeyValueArray } from "@/models/key-value";

export async function transferFrom(
  collectionAddress: string,
  tokenId: number,
  to: string,
  walletProvider: Eip1193Provider
) {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    connectedSigner
  );

  let tx;

  try {
    tx = await contract.transferFrom(
      await connectedSigner.getAddress(),
      to,
      tokenId
    );
    await tx.wait();
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }

  return tx;
}

export async function nameToTokenId(
  collectionAddress: string,
  name: string,
  walletProvider: Eip1193Provider
): Promise<bigint> {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    connectedSigner
  );

  return await contract.nameToTokenId(name);
}

export async function walletToTokenId(
  collectionAddress: string,
  name: string,
  walletProvider: Eip1193Provider
): Promise<bigint> {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    connectedSigner
  );

  return await contract.walletToTokenId(name);
}

export async function moderatorCount(
  collectionAddress: string,
  walletProvider: Eip1193Provider
): Promise<number> {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    connectedSigner
  );

  return await contract.moderatorCount();
}

export async function hasModeratorRole(
  collectionAddress: string,
  address: string,
  walletProvider: Eip1193Provider
): Promise<boolean> {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    connectedSigner
  );

  const moderatorRole = await contract.MODERATOR_ROLE();
  return await contract.hasRole(moderatorRole, address);
}

export async function createMintProposal(
  collectionAddress: string,
  signer: string,
  name: string,
  role: string,
  wallets: string[],
  images: Record<string, string>,
  attributes: Record<string, string>,
  cloneOf: number | null,
  nonce: string,
  signature: string,
  walletProvider: Eip1193Provider
) {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    connectedSigner
  );

  let tx;

  let cloneOfStruct = { isSet: false, value: 0 };
  if (cloneOf !== null) {
    cloneOfStruct = { isSet: true, value: cloneOf };
  }

  const imagesArray = mapToKeyValueArray(images);
  const attributesArray = mapToKeyValueArray(attributes);

  try {
    tx = await contract.createMintProposal(
      await connectedSigner.getAddress(),
      signer,
      role,
      name,
      wallets,
      imagesArray,
      attributesArray,
      cloneOfStruct,
      nonce,
      signature
    );

    await tx.wait();
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }

  return tx;
}

export async function mint(
  collectionAddress: string,
  signer: string,
  name: string,
  role: string,
  wallets: string[],
  images: Record<string, string>,
  attributes: Record<string, string>,
  cloneOf: number | null,
  nonce: string,
  signature: string,
  walletProvider: Eip1193Provider
) {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    connectedSigner
  );

  let tx;

  let cloneOfStruct = { isSet: false, value: 0 };
  if (cloneOf !== null) {
    cloneOfStruct = { isSet: true, value: cloneOf };
  }

  const imagesArray = mapToKeyValueArray(images);
  const attributesArray = mapToKeyValueArray(attributes);

  try {
    tx = await contract.mint(
      await connectedSigner.getAddress(),
      signer,
      role,
      name,
      wallets,
      imagesArray,
      attributesArray,
      cloneOfStruct,
      nonce,
      signature
    );

    await tx.wait();
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }

  return tx;
}

export async function updateAgentSigner(
  collectionAddress: string,
  tokenId: number,
  newSigner: string,
  nonce: string,
  signature: string,
  walletProvider: Eip1193Provider
) {
  const provider = new BrowserProvider(walletProvider);
  const frontSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    frontSigner
  );

  let tx;

  try {
    tx = await contract.updateAgentSigner(tokenId, newSigner, nonce, signature);

    await tx.wait();
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }

  return tx;
}

export async function addOrUpdateAgentImageFromProposal(
  collectionAddress: string,
  proposalIndex: number,
  walletProvider: Eip1193Provider
) {
  const provider = new BrowserProvider(walletProvider);
  const frontSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    frontSigner
  );

  let tx;

  try {
    tx = await contract.addOrUpdateAgentImageFromProposal(proposalIndex);

    await tx.wait();
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}

export async function refuseAddOrUpdateAgentImageProposal(
  collectionAddress: string,
  proposalIndex: number,
  walletProvider: Eip1193Provider
) {
  const provider = new BrowserProvider(walletProvider);
  const frontSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    frontSigner
  );

  let tx;

  try {
    tx = await contract.refuseAddOrUpdateAgentImageProposal(proposalIndex);

    await tx.wait();
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}

export async function mintFromProposal(
  collectionAddress: string,
  proposalIndex: number,
  walletProvider: Eip1193Provider
) {
  const provider = new BrowserProvider(walletProvider);
  const frontSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    frontSigner
  );

  let tx;

  try {
    tx = await contract.mintFromProposal(proposalIndex);

    await tx.wait();
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }

  return tx;
}

export async function refuseMintProposal(
  collectionAddress: string,
  proposalIndex: number,
  walletProvider: Eip1193Provider
) {
  const provider = new BrowserProvider(walletProvider);
  const frontSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    frontSigner
  );

  let tx;

  try {
    tx = await contract.refuseMintProposal(proposalIndex);

    await tx.wait();
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}

export async function getMintProposals(
  collectionAddress: string,
  walletProvider: Eip1193Provider
): Promise<MintProposal[]> {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    connectedSigner
  );

  try {
    const mintProposalsLength = await contract.mintProposalsLength();
    const promises = [];

    for (let i = 0; i < mintProposalsLength; i++) {
      promises.push(contract.mintProposal(i));
    }

    const mintProposals = await Promise.all(promises);

    return mintProposals.map((proposal) => ({
      id: proposal.id,
      collectionAddress: collectionAddress,
      name: proposal.name,
      role: proposal.role,
      owner: proposal.owner,
      signer: proposal.signer,
      mission: proposal.mission,
      creator: proposal.creator,
      wallets: proposal.wallets,
      images: proposal.images,
    }));
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}

export async function getAddOrUpdateAgentImageProposals(
  collectionAddress: string,
  walletProvider: Eip1193Provider
): Promise<AddOrUpdateImageProposal[]> {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    connectedSigner
  );

  try {
    const addOrUpdateAgentImageProposalsLength =
      await contract.addOrUpdateAgentImageProposalsLength();
    const promises = [];

    for (let i = 0; i < addOrUpdateAgentImageProposalsLength; i++) {
      promises.push(contract.addOrUpdateAgentImageProposal(i));
    }

    const addOrUpdateAgentImageProposals = await Promise.all(promises);

    return addOrUpdateAgentImageProposals.map((proposal) => ({
      id: proposal.id,
      collectionAddress: collectionAddress,
      tokenId: proposal.tokenId,
      image: {
        key: proposal.image.key,
        value: proposal.image.value,
      },
    }));
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}

export async function addOrUpdateAgentAttribute(
  collectionAddress: string,
  tokenId: number,
  attributeKey: string,
  attributeValue: string,
  walletProvider: Eip1193Provider
) {
  const provider = new BrowserProvider(walletProvider);
  const frontSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    collectionAddress,
    AgentCollectionABI,
    frontSigner
  );

  let tx;

  try {
    tx = await contract.addOrUpdateAgentAttribute(
      tokenId,
      attributeKey,
      attributeValue
    );

    await tx.wait();
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}
