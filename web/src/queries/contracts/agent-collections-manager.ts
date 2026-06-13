import { Address } from "@reown/appkit-adapter-ethers";
import { BrowserProvider, Eip1193Provider, ethers } from "ethers";
import { AgentCollectionsManagerABI } from "@/abi/agent-collections-manager";

export async function moderatorOf(
  walletProvider: Eip1193Provider
): Promise<Address[]> {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    import.meta.env.VITE_AGENT_COLLECTIONS_MANAGER_ADDRESS,
    AgentCollectionsManagerABI,
    connectedSigner
  );

  try {
    return await contract.moderatorOf(await connectedSigner.getAddress());
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionsManagerABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}

export async function collectionIdByName(
  name: string,
  walletProvider: Eip1193Provider
): Promise<bigint> {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    import.meta.env.VITE_AGENT_COLLECTIONS_MANAGER_ADDRESS,
    AgentCollectionsManagerABI,
    connectedSigner
  );

  try {
    return await contract.collectionIdByName(name);
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionsManagerABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}

export async function collectionIdBySymbol(
  symbol: string,
  walletProvider: Eip1193Provider
): Promise<bigint> {
  const provider = new BrowserProvider(walletProvider);
  const connectedSigner = await provider.getSigner();

  const contract = new ethers.Contract(
    import.meta.env.VITE_AGENT_COLLECTIONS_MANAGER_ADDRESS,
    AgentCollectionsManagerABI,
    connectedSigner
  );

  try {
    return await contract.collectionIdBySymbol(symbol);
  } catch (error) {
    const iface = new ethers.Interface(AgentCollectionsManagerABI);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    console.error(error, iface.parseError((error as any).data));

    throw error;
  }
}
