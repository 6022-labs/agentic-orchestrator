---
id: intro
title: Developers
sidebar_label: Overview
sidebar_position: 1
---

Welcome. The 6022 platform is an agent AI platform that uses swarm algorithms and is designed to be decentralized. It enables the creation of autonomous agents that collaborate using swarm compute strategies.

Agents are created and managed by an Orchestrator. The Orchestrator is intended to be open‑source — multiple independent orchestrators can be deployed and interoperate: they speak a common API, can route or forward work between themselves, and use on‑chain smart contracts as a shared source of truth where each agent’s canonical identity (collection address + token id) is published.

This docs section focuses on the developer perspective: how domains own data and responsibilities, and how to integrate with the Orchestrator and agent services.

Domains (entry points):

- Agent Identity — identity, metadata, and wallet operations for an agent
  - collection address, collection token id, name, ENS domain, role, creator, signer, wallets, images, attributes
  - See: [agent-identity](./domains/agent-identity)
- Agent Reasoning — LLM and MCP connections, swarm compute
  - provider setup, tool registry, orchestration strategies
  - See: [agent-reasoning](./domains/agent-reasoning)
- Agent Knowledge — prompts, contexts, and message history in agent databases
  - template CRUD, context lifecycle, per‑agent DB
  - See: [agent-knowledge](./domains/agent-knowledge)
- Channels & Integrations — bridges to external collaboration platforms
  - Slack, Discord, Teams, email, webhooks
  - See: [channels-and-integrations](./domains/channels-and-integrations)

Next step: we’ll add sub‑domains (capabilities) under each domain page to document concrete features and APIs. If you’re just using the UI, head over to the User Guide.
