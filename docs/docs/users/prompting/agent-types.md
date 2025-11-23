---
title: Agent Types
sidebar_position: 2
---

# Agent Types

> Part of [Agent Prompting](./intro.md) documentation

## Introduction

The 6022 framework defines four distinct agent types, each with specific roles, capabilities, and interaction rules. Understanding these types is essential for designing effective multi-agent systems.

---

## Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         AGENT TYPES                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐            │
│    │  FACILITATOR │    │    EXPERT    │    │    HUMAN     │            │
│    │              │    │              │    │              │            │
│    │  Orchestrates│    │  Specialized │    │  Real person │            │
│    │  conversation│    │  knowledge   │    │  in the loop │            │
│    └──────┬───────┘    └──────┬───────┘    └──────┬───────┘            │
│           │                   │                   │                    │
│           │                   │                   │                    │
│           │            ┌──────┴───────┐           │                    │
│           │            │ DIGITAL TWIN │           │                    │
│           │            │              │◄──────────┘                    │
│           │            │  AI clone of │     FORK                       │
│           │            │  a human     │                                │
│           │            └──────────────┘                                │
│           │                                                            │
│           ▼                                                            │
│    Calls experts by ENS name                                           │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Agent Facilitator

The **Facilitator** is the conversation orchestrator. It receives initial requests, coordinates expert agents, and produces the final response.

### Characteristics

| Property                | Value                                 |
| ----------------------- | ------------------------------------- |
| **Role**                | Orchestrate multi-agent conversations |
| **Prompts**             | System, Default, Conclude             |
| **Context Primary Key** | Channel-based (conversation itself)   |
| **Memory**              | Full conversation history             |
| **Can call**            | All agents (not self)                 |
| **Triggered by**        | Initial request or timeout            |

### Key Behaviors

1. **Receives initial request** — First to process incoming messages
2. **Calls experts** — Mentions agents by full ENS name to get information
3. **Synthesizes information** — Aggregates responses from experts
4. **Manages conversation flow** — Decides when to continue or conclude
5. **Produces final response** — Generates conclusion for the user

### Prompt Structure

```
┌────────────────────────────────────────┐
│           SYSTEM PROMPT                │
│  Identity, mission, participants       │
├────────────────────────────────────────┤
│           DEFAULT PROMPT               │
│  How to manage conversation            │
│  When to call experts                  │
│  How to end conversation               │
├────────────────────────────────────────┤
│          CONCLUDE PROMPT               │
│  How to write final response           │
│  Format requirements                   │
└────────────────────────────────────────┘
```

---

## Agent Expert

**Expert** agents hold specialized knowledge in a specific domain. They respond when called by the facilitator or other agents.

### Characteristics

| Property                | Value                                                |
| ----------------------- | ---------------------------------------------------- |
| **Role**                | Provide specialized knowledge                        |
| **Prompts**             | System, Context Found, Context Not Found             |
| **Context Primary Key** | Domain-specific (customer ID, contract number, etc.) |
| **Memory**              | Only messages where mentioned                        |
| **Can call**            | All agents (not self)                                |
| **Triggered by**        | Being mentioned by ENS name                          |

### Key Behaviors

1. **Triggered by mention** — Only responds when explicitly named
2. **Threading phase** — Identifies context primary key on first call
3. **Answers with context** — Uses memory associated with the key
4. **May ask for help** — Can call other agents for additional info
5. **Cannot self-call** — Prevented from mentioning own name

### Threading Behavior

```
1. First mention: Agent enters THREADING state
2. Key detection: Analyzes message for context primary key
3. Key lock: Once found, key is fixed for conversation duration
4. Memory load: All history for that key becomes available
```

---

## Agent Human

**Human** agents represent real people participating in the conversation via Slack or Teams.

### Characteristics

| Property                | Value                                       |
| ----------------------- | ------------------------------------------- |
| **Role**                | Human in the loop                           |
| **Prompts**             | None (uses Slack/Teams interface)           |
| **Context Primary Key** | Unordered list of participant ETH addresses |
| **Memory**              | Full exchange history                       |
| **Can call**            | Any agent by ENS name                       |
| **Triggered by**        | Manual message from Slack/Teams             |

### Key Behaviors

1. **"Biological LLM"** — Human prompted by the system for responses
2. **No threading** — Key is known upon channel entry
3. **Database writing** — Exchanges stored for potential cloning
4. **Traceability** — Identified by full ENS in conversation
5. **Superpower** — Can restart a closed conversation

---

## Agent Digital Twin

**Digital Twin** agents are AI clones of human agents, created via a FORK operation.

### Characteristics

| Property                | Value                                                |
| ----------------------- | ---------------------------------------------------- |
| **Role**                | AI substitute for a human                            |
| **Prompts**             | System + Default (both written by human counterpart) |
| **Context Primary Key** | Participant addresses (automatic, like human)        |
| **Threading**           | None — context primary key determined automatically  |
| **Memory**              | Copy of human's full history                         |
| **Can call**            | All agents (not self)                                |
| **Triggered by**        | Being mentioned by ENS name                          |

### Key Behaviors

1. **Origin** — Created by duplicating human's conversation history
2. **Unique identity** — Own NFT, ENS, and wallet
3. **System prompt** — Written by the human counterpart to define identity
4. **Default prompt** — Equivalent to expert's "Context Found", enables swarm participation
5. **No threading phase** — Context primary key is automatic (interlocutor set)
6. **Configurable LLM** — Owner chooses the underlying model
7. **Human link** — Keeps original human in whitelist for arbitration
8. **Priority rule** — If human responds, human's answer prevails
9. **Interaction limits** — Cannot self-call

---

## Interaction Rules

### Who Can Call Whom

| Caller           | Can Call   | Cannot Call |
| ---------------- | ---------- | ----------- |
| **Facilitator**  | All agents | Self        |
| **Expert**       | All agents | Self        |
| **Human**        | All agents | —           |
| **Digital Twin** | All agents | Self        |

### How to Call an Agent

Agents are triggered by mentioning their **full ENS domain name**:

```
@expert-alpha.137.protocol6022.eth, what is the status
of customer C45653?
```

**Important:** Short names do not trigger agents. Only the full ENS works.

---

## Memory Management by Type

| Agent Type       | Memory Scope              | Context Primary Key   |
| ---------------- | ------------------------- | --------------------- |
| **Facilitator**  | Full conversation         | Conversation itself   |
| **Expert**       | Only @-mentioned messages | Domain-specific       |
| **Human**        | Full exchange history     | Participant addresses |
| **Digital Twin** | Inherited + new exchanges | Same as human         |

---

## Summary

| Type             | Purpose               | Prompts                         | Key Feature                   |
| ---------------- | --------------------- | ------------------------------- | ----------------------------- |
| **Facilitator**  | Orchestration         | System, Default, Conclude       | Sees full conversation        |
| **Expert**       | Specialized knowledge | System, Context Found/Not Found | Context primary key isolation |
| **Human**        | Human in the loop     | None                            | Can restart conversations     |
| **Digital Twin** | AI substitute         | System + Default                | Defers to original human      |

---

## Related Documentation

- [Conversation Flow](./conversation-flow.md) — How agents interact over time
- [Facilitator Prompts](./facilitator-prompts.md) — Configuring facilitator behavior
- [Expert Prompts](./expert-prompts.md) — Configuring expert behavior
