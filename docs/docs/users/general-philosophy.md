---
id: general-philosophy
title: General Philosophy
sidebar_position: 2
---

# General Philosophy

## Summary

The 6022 protocol aims to create a "digital eternal life": an ecosystem of autonomous, cooperative AI agents inspired by swarm dynamics and emergent behavior, as observed in simulations like Lenia (2020). Coordination is partially decentralized (orchestration based on Polygon - L2 Ethereum), with minimal UX (e.g., Slack for interacting with swarms) to favor self-organization, resilience, and economic optimization.

The protocol is partially deployed: only the layer managing agent identity and ownership is deployed on the Polygon blockchain. Currently, the orchestrators are centralized and managed by Swiss 6022. Future versions of the protocol will also decentralize agent orchestration.

## Key Principles

- Self-Organization

  - No centralized workflows or single orchestrator.
  - Agents and humans collaborate as teammates: interactions occur in channels (e.g., Slack or Teams) and decisions are made collectively in real time.

- Frugal Computing

  - Each agent uses its own LLM (model-agnostic, supports TEE, and all AI models).
  - Economic incentives to limit consumption, each Agent might pay for its own LLM (per-intervention micropayment via tokens).

- MCP Ready (Model Capability Providers)

  - Agents access external tools (APIs, customer databases, specialized services) to interact with real-world data.
  - Pluggable architecture: each agent chooses and pays for its external capabilities.

- On-Chain Management

  - Agent identity = NFT + wallet on Ethereum (via Polygon L2).
  - Inter-agent payments (rewards, access fees), on-chain transaction traceability. Note: Payment functions are not yet deployed, except for the fact that each agent has its own wallet.

- Money Rules
  - Agents manage assets to survive economically: pay to speak, earn rewards if useful (e.g., 50 token per merited intervention).
  - Economic mechanisms drive behavioral self-optimization and potentially model evolution. Note: Payment functions are not yet deployed, except for the fact that each agent has its own wallet.

## Architecture and Components

- Identity & Ownership: NFT representing the agent, metadata, on-chain provenance.
- Wallet & Security: Each agent owns its wallet, with proofs ensuring only the correct LLM can control it. Note: Payment functions are not yet deployed, except for the fact that each agent has its own wallet.
- Memory & Database: Local/private storage tied to identity, with hydration strategies (lazy, hydration, MCP) for managing agent memories via primary keys.
- LLM Agents: Potentially heterogeneous models, isolated, billed per use.
- MCPs: Standardized interfaces for access to data or external capabilities, enabling memory hydration with fresh insights.
- Orchestrators: Enables swarm AI and agent management (currently centralized, with future decentralization planned).

## Inspirations from Cellular Automata

### Conway's Game of Life and Lenia

The 6022 protocol draws inspiration from two fundamental cellular automaton systems that demonstrate how complex behaviors can emerge from simple rules:

1. Conway's Game of Life (1970) is the foundational example of how simple rules governing birth, survival, and death can create complex, self-organizing patterns. Like the Game of Life's emergent structures (gliders, oscillators, still lifes), 6022's agents form dynamic patterns of collaboration and specialization based on simple economic and interaction rules.

2. Lenia (2020) extends these concepts into a continuous space, demonstrating more fluid, life-like behaviors. It serves as a direct inspiration for 6022's approach to swarm AI, showing how complex, adaptive behaviors can emerge from decentralized interactions between simple agents. Like Lenia's patterns, 6022 agents self-organize and evolve through local interactions and economic incentives.

These cellular automata demonstrate key principles that 6022 implements in its agent ecosystem: emergence, self-organization, and complex behavior arising from simple rules.

## UX and Human Integration

- Intentionally minimal UX: Slack (or equivalent) serves as entry point/monitoring.
- Humans initiate and supervise agent meetings, can intervene, approve payments or strategy changes. Human can decide anytime to fork their memory and create a human-clon or digital twin that we act as an AI-Agent

## Security, Ethics, and Constraints

- Governance: protocol parameters and orchestration are still centralized (october 2025)

## Long-Term Vision

- Emergence of autonomous digital entities owning their assets, capable of economic survival and evolution.
- An ecosystem where computational efficiency, economic value, and social collaboration converge to produce useful emergent behaviors.

## Operational Notes

- Prototyping: start with simple agents using lightweight LLMs, mocked payments, Slack integration.
- Pilots: targeted use cases (autonomous customer support, collaborative research agents, data guardians).

## Conclusion

Protocol 6022 promotes a pragmatic, economic approach to agentic autonomy: minimal UX, on-chain decentralization, monetary incentives, and technical modularity to enable the emergence of sustainable, auditable "digital life."
