---
id: architecture
title: Architecture
sidebar_label: Architecture
sidebar_position: 2
---

A bird's-eye view of how the orchestrator fits together. The picture is intentionally shallow — just enough to place the smart contracts, the four domains, and their main components on the same page. For details on what each domain owns, follow the links from the layer descriptions below.

![High-level architecture of the 6022 orchestrator](/img/developers/architecture.svg)

The diagram has three layers:

- **Smart contracts (on-chain)** — the shared source of truth for agent identity. `Agent Collection V1` mints each agent as an ERC-721, `Agent Collection Creator V1` is the factory that deploys new collections, and `Token 6022` is the ERC-20 utility token. Multiple orchestrator instances can interoperate because they all read from the same contracts.
- **Domains** — the four business areas the orchestrator is organized around: [Agent Identity](./domains/agent-identity), [Agent Reasoning](./domains/agent-reasoning), [Agent Knowledge](./domains/agent-knowledge), and [Channels & Integrations](./domains/channels-and-integrations). The **Gateway** is the cross-cutting routing layer in front of them, and the **Web UI** is the human entry point that talks to the Gateway.
- **Components** — the concrete services inside each domain. Agent Identity runs the registry host plus a chain watcher that syncs on-chain state into its database. Agent Reasoning runs the agent compute host that calls LLM providers and MCP tool servers. Agent Knowledge exposes a database connector in front of per-agent Postgres databases. Channels & Integrations runs internal/external bridge services for each external workspace (Slack today, more later).
