// Display names + public RPCs for the EVM networks the console knows. Used to
// label the network select and to look up ERC-20 metadata while typing; the
// set of selectable networks comes from the runtime's facilitator config.
export type EvmNetwork = {
  // CAIP-2 identifier, e.g. "eip155:80002".
  id: string
  name: string
  // Public read-only RPC used for ERC-20 metadata lookups while typing.
  rpcUrl: string
  // Testnet chains resolve ENS via the Sepolia app, mainnet via the live one.
  testnet: boolean
}

export const EVM_NETWORKS: EvmNetwork[] = [
  {
    id: "eip155:1",
    name: "Ethereum",
    rpcUrl: "https://eth.llamarpc.com",
    testnet: false,
  },
  {
    id: "eip155:137",
    name: "Polygon",
    rpcUrl: "https://polygon-rpc.com",
    testnet: false,
  },
  {
    id: "eip155:8453",
    name: "Base",
    rpcUrl: "https://mainnet.base.org",
    testnet: false,
  },
  {
    id: "eip155:42161",
    name: "Arbitrum One",
    rpcUrl: "https://arb1.arbitrum.io/rpc",
    testnet: false,
  },
  {
    id: "eip155:11155111",
    name: "Ethereum Sepolia",
    rpcUrl: "https://ethereum-sepolia-rpc.publicnode.com",
    testnet: true,
  },
  {
    id: "eip155:80002",
    name: "Polygon Amoy",
    rpcUrl: "https://rpc-amoy.polygon.technology",
    testnet: true,
  },
  {
    id: "eip155:84532",
    name: "Base Sepolia",
    rpcUrl: "https://sepolia.base.org",
    testnet: true,
  },
  {
    id: "eip155:421614",
    name: "Arbitrum Sepolia",
    rpcUrl: "https://sepolia-rollup.arbitrum.io/rpc",
    testnet: true,
  },
  {
    id: "eip155:5115",
    name: "Citrea Testnet",
    rpcUrl: "https://rpc.testnet.citrea.xyz",
    testnet: true,
  },
  {
    id: "eip155:99999",
    name: "ADI Testnet",
    rpcUrl: "https://rpc.ab.testnet.adifoundation.ai",
    testnet: true,
  },
]

export function findEvmNetwork(id: string): EvmNetwork | undefined {
  return EVM_NETWORKS.find((network) => network.id === id)
}
