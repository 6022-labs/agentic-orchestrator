---
title: Agent Types
sidebar_position: 2
---

# Agent Types

Four types of agents. Each has a specific job and a specific memory model — picking the right type is the first decision you make when designing a swarm.

| Type | Purpose | Prompts | Triggered by | Memory |
|---|---|---|---|---|
| **Facilitator** | Orchestrate the conversation, call experts, produce final response | System, Default, Conclude | Initial request or timeout | Current conversation only |
| **Expert** | Provide specialised domain knowledge | System, Context Found, Context Not Found | Being @-mentioned by full ENS | Cross-conversation, keyed by a domain-specific context primary key |
| **Human** | Real person in the loop | None (types in Slack) | Manual message from Slack/Teams | Full history, keyed automatically by interlocutor set |
| **Digital Twin** | AI clone of a human, built by FORK | System, Default | Being @-mentioned | Inherited from the human + new exchanges, keyed like a human |

## Facilitator

The orchestrator. Sees every message in the thread, calls experts by full ENS name, decides when to conclude and produces the final user-facing response. Has no memory across conversations — each new thread starts fresh.

See [Facilitator Prompts](./facilitator-prompts.md) for the three prompts (System, Default, Conclude) and how to configure them.

## Expert

A domain specialist. Only sees messages where it's @-mentioned, but maintains memory across every conversation that references the same key. The key is domain-specific — a customer ID, a contract number, an article reference.

On first mention in a thread, the expert runs a **threading phase**: it analyses the message, identifies the key, locks it for the rest of the conversation, and loads the matching memory. See [Expert Prompts](./expert-prompts.md) for the three prompts (System, Context Not Found, Context Found) and [Threading Phase](../memory/threading-phase.md) for what happens under the hood.

## Human

A real person. No prompts — humans type directly into Slack. Their context primary key is automatically the set of participant wallet addresses in the conversation, so the "same human, different group" scenario produces a different memory. Humans are the only agents that can **restart a concluded conversation** by posting a new message.

See [Human & Digital Twin](./human-digital-twin.md) for setup (wallet verification, Slack pairing).

## Digital Twin

An AI clone of a human, created by a FORK operation. Inherits the human's full conversation history at the moment of the fork, then accumulates its own. Has its own NFT, ENS and wallet. Skips the threading phase entirely because its key is automatic (interlocutor set, like a human).

When the twin is uncertain, it can call the original human for arbitration — if the human responds, the human's answer takes absolute priority over the twin's.

See [Human & Digital Twin](./human-digital-twin.md) for the full story.

## Interaction rules

Every agent can call every other agent by full ENS name. The only restriction: **no agent can call itself** (this prevents infinite loops). Humans have one extra power that no other type has: they can restart a conversation after the facilitator has concluded it.

Short names don't trigger anything. Only the full ENS domain (`name.collection.chainId.protocol6022.eth`) works.
