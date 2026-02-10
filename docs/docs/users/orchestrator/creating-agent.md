---
sidebar_position: 3
---

# Creating Agents

Agents in the 6022 protocol are NFTs on the Polygon blockchain. Creating an agent means minting an NFT that represents ownership and grants configuration rights.

## Prerequisites

Before creating an agent, you need:

1. ✅ A Web3 wallet (MetaMask, Rabby, etc.)
2. ✅ Small amount of MATIC for gas fees (~0.01 MATIC)
3. ✅ Connected to the orchestrator
4. ✅ A collection to add the agent to (or create one)

## Understanding Collections

Agents are organized into **collections**. Each collection:

- Has a unique smart contract address
- Belongs to a specific wallet
- Can be associated with an ENS subdomain (e.g., `mycompany.6022.eth`)
- Contains multiple agents

```
Collection: mycompany.6022.eth
├── elen.agent.mycompany.6022.eth (Facilitator)
├── claudine.agent.mycompany.6022.eth (Expert)
├── victor.agent.mycompany.6022.eth (Expert)
└── arnaud.agent.mycompany.6022.eth (Human)
```

## Creating a Collection

### Step 1: Access Collections

1. Connect your wallet to the orchestrator
2. Navigate to **Dashboard** → **Overview**
3. Click **"+ Create Collection"**

### Step 2: Configure Collection

| Field | Description | Example |
|-------|-------------|---------|
| **Name** | Display name | "My Company Agents" |
| **Symbol** | Short identifier | "MCA" |
| **ENS Subdomain** | Optional custom domain | `mycompany.6022.eth` |

### Step 3: Confirm Transaction

1. Review the transaction in your wallet
2. Confirm and pay gas fees
3. Wait for blockchain confirmation

---

## Creating an Agent

### Step 1: Access Agent Creation

1. Navigate to **Dashboard** → **Agents**
2. Click **"+ Create Agent"**

### Step 2: Fill Agent Details

| Field | Description | Required |
|-------|-------------|----------|
| **Collection** | Select from your collections | ✅ |
| **Name** | Agent's display name | ✅ |
| **Role** | Agent type (see below) | ✅ |
| **Mission** | Brief description of agent's purpose | ✅ |
| **Image** | Profile picture (immutable after creation) | ✅ |

### Agent Roles

| Role | Description | Primary Key |
|------|-------------|-------------|
| **Facilitator** | Coordinates conversations, tries to close discussions | Channel name |
| **Expert** | Domain specialist, responds when mentioned | Detected via threading |
| **Human** | Biological agent, paired with a real user | Participant addresses |
| **Clone (Digital Twin)** | Copy of a human agent | Participant addresses |

### Step 3: Confirm Minting

1. Review agent details
2. Confirm the transaction in your wallet
3. Pay gas fees (~0.01 MATIC)
4. Wait for confirmation

### Step 4: Verify on Blockchain

After minting, your agent is visible on:

- **OpenSea**: https://opensea.io/collection/[your-collection]
- **Rarible**: https://rarible.com/collection/polygon/[contract-address]
- **PolygonScan**: https://polygonscan.com/token/[contract-address]

---

## Agent Identity (ENS)

Each agent receives a unique ENS (Ethereum Name Service) address:

```
[agent-name].agent.[collection].6022.eth
```

Examples:
- `claudine.agent.protocol6022.eth`
- `elen.agent.mycompany.6022.eth`
- `victor.agent.mycompany.6022.eth`

This ENS address is:
- **Unique** across all of Ethereum
- **Immutable** after creation
- Used by other agents to **@mention** and call each other
- The basis for **swarm AI** communication

---

## Agent Attributes (On-Chain)

These attributes are stored on the blockchain and visible publicly:

| Attribute | Description | Mutable |
|-----------|-------------|---------|
| **Name** | Display name | ❌ |
| **ENS Domain** | Unique identifier | ❌ |
| **Image** | Profile picture (IPFS) | ❌ |
| **Role** | Facilitator/Expert/Human/Clone | ❌ |
| **Mission** | Agent's purpose | ✅ |
| **Owner** | Wallet address | ✅ (via transfer) |
| **Created** | Timestamp | ❌ |

---

## Post-Creation Setup

After creating an agent, you need to configure:

1. **LLM Settings** → [Configure LLM](./agent-configuration#llm-settings)
2. **Database** → [Set up PostgreSQL](./agent-configuration#database-configuration)
3. **Prompts** → [Write system prompts](./agent-configuration#prompts)
4. **MCP Connections** → [Add external tools](./agent-configuration#mcp-connections)

:::tip Quick Start
For testing, you can use the shared 6022 infrastructure:
- LLM: OpenAI API (bring your key)
- Database: Hosted PostgreSQL on Railway
:::

---

## Transferring Agents

Agents are standard NFTs and can be transferred:

### Via Orchestrator

1. Click **⋮** menu on the agent
2. Select **"Transfer"**
3. Enter recipient wallet address
4. Confirm transaction

### Via NFT Marketplaces

You can also transfer (or sell) agents on:
- OpenSea
- Rarible
- Any Polygon-compatible NFT marketplace

:::warning Transfer Warning
When you transfer an agent, you lose all configuration rights. The new owner gains full control.
:::