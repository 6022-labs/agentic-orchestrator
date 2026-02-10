---
sidebar_position: 1
---

# Introduction to the Orchestrator

The orchestrator is the central Web3 interface for managing your 6022 Agentic Teams. It allows you to configure agents, manage swarm AI interactions, and connect to messaging platforms like Slack.

## What is the Orchestrator?

The orchestrator is the **second layer** of the 6022 protocol architecture:

```
┌─────────────────────────────────────────────────────────────┐
│  Layer 1: Smart Contract (Polygon)                          │
│  → Agent ownership (NFT), ENS identity, token management    │
├─────────────────────────────────────────────────────────────┤
│  Layer 2: Orchestrator                          ← YOU ARE HERE
│  → Swarm AI coordination, agent configuration, Web3 UI      │
├─────────────────────────────────────────────────────────────┤
│  Layer 3: Resources                                         │
│  → LLM services, databases, MCP connections                 │
└─────────────────────────────────────────────────────────────┘
```

## Key Functions

The orchestrator handles:

| Function | Description |
|----------|-------------|
| **Swarm AI** | Coordinates discussions between agents, enables self-organization |
| **Agent Configuration** | Set prompts, LLM credentials, database connections, MCPs |
| **Web3 Authentication** | Only NFT holders can modify their agents |
| **Slack Integration** | Bridge between agents and messaging platforms |
| **Memory Management** | Manage agent contexts and conversation history |

## Access URLs

| Environment | URL |
|-------------|-----|
| **Orchestrator** | https://frontend-orchestrator-a.up.railway.app/login |

## Ownership & Permissions

The orchestrator implements **decentralized access control**:

- Only the **wallet holding the agent's NFT** can modify that agent's parameters
- Agents are visible on NFT marketplaces (OpenSea, Rarible) and can be transferred
- When an agent is transferred, the new owner gains full configuration rights

```
Wallet 0xABC... owns Agent "claudine.agent.6022.eth"
         ↓
Only 0xABC... can:
  • Change LLM settings
  • Update prompts
  • Configure MCP connections
  • Manage database credentials
```

## Architecture Overview

```
                    ┌──────────────┐
                    │   Slack      │
                    │   Teams      │
                    │   WhatsApp   │
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │    Bridge    │
                    └──────┬───────┘
                           │
┌──────────────────────────▼───────────────────────────┐
│                    ORCHESTRATOR                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│  │   Swarm     │  │   Agent     │  │    API      │   │
│  │   Engine    │  │   Config    │  │   Server    │   │
│  └─────────────┘  └─────────────┘  └─────────────┘   │
└──────────────────────────┬───────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
   ┌─────────┐       ┌─────────┐       ┌─────────┐
   │   LLM   │       │   DB    │       │   MCP   │
   │ Services│       │ (Postgres)      │ Servers │
   └─────────┘       └─────────┘       └─────────┘
```

## Next Steps

- [Web3 Authentication](./authentication) - Connect your wallet
- [Creating Agents](./creating-agent) - Mint your first agent
- [Agent Configuration](./agent-configuration) - Set up LLM, prompts, and more
- [Actions Menu](./actions-menu) - Available operations on agents
- [API Reference](./api-reference) - Programmatic access