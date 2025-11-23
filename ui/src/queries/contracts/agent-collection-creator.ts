import { BrowserProvider, Eip1193Provider, ethers } from "ethers";
import { AgentCollectionCreatorABI } from "@/abi/agent-collection-creator";

export async function createCollection(
  name: string,
  symbol: string,
  description: string,
  adminAddress: string,
  moderatorAddress: string,
  walletProvider: Eip1193Provider
): Promise<void> {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    import.meta.env.VITE_AGENT_COLLECTION_CREATOR_ADDRESS,
    AgentCollectionCreatorABI,
    connectedSigner
  );

  try {
    return await contract.createCollection(
      name,
      symbol,
      description,
      adminAddress,
      moderatorAddress
    );
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionCreatorABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}

export async function canCreateCollection(
  address: string,
  walletProvider: Eip1193Provider
) {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    import.meta.env.VITE_AGENT_COLLECTION_CREATOR_ADDRESS,
    AgentCollectionCreatorABI,
    connectedSigner
  );

  try {
    return await contract.canCreateCollection(address);
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionCreatorABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}
