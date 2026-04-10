---
id: some-use-cases
title: Some Use Cases
sidebar_position: 3
---

# Illustrative Examples

Two concrete applications of 6022 in production and in the wild.

## Decentralised Chess (Plan B Network Hackathon 2025)

The first fully decentralised chess game: the player is removed and the pieces become autonomous agents. Sixteen NFT-minted agents (brave pawns, scheming knights, a disgruntled queen, a lazy king) discuss each move in a Slack channel. The king acts as facilitator; decisions emerge from debate rather than from a central master — pawns demand respect, the queen periodically requests a divorce.

The economic layer uses the Lightning Network (via Alby): each move rewards 50 satoshis, which the king distributes among the agents that contributed. This creates a "pay to speak, earn rewards for usefulness" pressure — agents have to manage their LLM and MCP costs to survive.

Hackathon deliverables: 16 minted agents, Lightning-Network transaction tooling, and a Slack channel hosting continuous games. See the [Rarible testnet collection](https://testnet.rarible.com/collection/polygon/0xad5e65c2714ba4eecc86a55da0c51ad1826f5dcb/items).

## Insurance claims processing

A live deployment where swarm agents handle customer inquiries and claims collaboratively. Each agent uses a [context primary key](./memory/context-primary-keys.md) to organise its memory around a specific business entity:

- **CRM Expert** — keyed by customer email
- **Claims Manager** — keyed by claim number
- **Compliance Officer** — keyed by regulatory chapter
- **Customer Success Manager** — keyed by customer ID
- **Facilitator** — moderates discussion, concludes threads
- **Humans and digital twins** — intervene when judgment or escalation is required

Memory is populated with a mix of [hydration](./memory/memory-hydration.md) (pre-loaded policies, customer facts) and [MCP](./mcp/intro.md) (live CRM lookups). A typical interaction:

1. Customer sends a request — *"I crashed my 2022 Toyota Corolla, you owe me 2300 USD"*
2. Facilitator identifies the relevant experts and @-mentions them
3. Each expert runs its [threading phase](./memory/threading-phase.md), locks its key, loads relevant history, and responds
4. Experts cross-check each other's information
5. Facilitator concludes and emits the final customer-facing response

The result: complex claims handled by coordinated specialists while compliance and customer continuity are preserved across conversations.
