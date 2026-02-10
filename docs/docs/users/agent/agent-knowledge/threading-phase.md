---
title: Threading Phase
sidebar_position: 5
---

# Threading Phase

## Introduction

The **threading phase** is the critical process during which an expert agent identifies which context primary key applies to a conversation when it is first mentioned.

This document explains the importance of context primary keys during threading and how agents identify them.

---

## Threading Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                        THREADING PHASE                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. Agent is mentioned          2. Agent enters              3. Agent searches    │
│     in conversation                THREADING state              for context key   │
│                                                                     │
│     "@claudine can you           Context = "THREADING"         "C45653" found     │
│      help with C45653?"                                        in message         │
│            │                            │                            │            │
│            ▼                            ▼                            ▼            │
│     ┌───────────┐               ┌───────────────┐            ┌───────────────┐   │
│     │  Mention  │──────────────▶│  THREADING    │───────────▶│ Context Found │   │
│     │  Detected │               │    State      │            │   = C45653    │   │
│     └───────────┘               └───────────────┘            └───────────────┘   │
│                                        │                            │            │
│                                        │ Key not found              │            │
│                                        ▼                            ▼            │
│                                 ┌───────────────┐            ┌───────────────┐   │
│                                 │ Ask for help  │            │ Load memory   │   │
│                                 │ from others   │            │ for C45653    │   │
│                                 └───────────────┘            └───────────────┘   │
│                                                                     │            │
│                                                                     ▼            │
│                                                              ┌───────────────┐   │
│                                                              │ Key locked    │   │
│                                                              │ for session   │   │
│                                                              └───────────────┘   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## How Threading Works

When an expert agent is first mentioned in a conversation, it triggers the threading process:

1. **Agent is mentioned** — Someone calls the agent to help with a specific issue
2. **Agent analyzes the message** — Looks for clues about which context primary key applies
3. **Key is identified** — Extracts or determines the correct context primary key
4. **Memory is loaded** — All associated memory for that key becomes available
5. **Key is locked** — Remains fixed for the duration of the conversation

---

## Key Identification Methods

Agents identify context primary keys through:

### 1. Explicit Mention

The key appears directly in the message:

```
"Can you help with customer C45653?"
→ Identifies: C45653
```

### 2. Pattern Matching

The key follows a known format (regex patterns, email format, UUID, etc.):

```
"Example Customer (customer@example.com) needs help"
→ Matches email pattern
→ Identifies: customer@example.com
```

### 3. MCP Lookup

The agent queries external systems to find the correct key:

```
"The customer from the automotive workshop"
→ Query CRM MCP
→ Identifies: CUST-45653
```

---

## What Happens After Threading

Once the context primary key is identified and locked:

- The agent loads all associated memory (from memory strategies: lazy, hydration, MCP)
- The agent engages in meaningful conversation with proper context
- The key remains locked for the duration of the conversation
- Future messages use the loaded memory to provide informed responses

---

## Configuration

Detailed configuration of key detection, prompt design, and MCP integration is covered in the **[Agent Prompting](../agent-prompting/intro.md)** section.

---

## Summary

The threading phase is critical for:

- **Correct context selection** — Identifying which context primary key applies
- **Memory isolation** — Loading the right conversation history
- **Conversation coherence** — Ensuring agents have proper context for meaningful responses

Without proper threading, agents cannot distinguish between different customers, contracts, tickets, or other business entities they manage.

---

## Related Documentation

- [Context Primary Keys](./context-primary-keys.md) - Creating and managing keys
- [Agent Prompting](../agent-prompting/intro.md) - Configuring key detection and prompts
- [Memory Strategies](./memory-strategies.md) - How memory is populated for a locked context
