---
title: Introduction
sidebar_position: 1
---

# Memory

> Documentation for the 6022 Agentic Teams project

## Overview

Memory encompasses everything an agent "knows" and "remembers" across conversations. In the 6022 framework, this knowledge is structured around **context primary keys** that organize memory into distinct contexts, **strategies** that determine how memory is populated, and **APIs** that enable programmatic management.

---

## Documentation Structure

### Core Concepts

| #   | Topic                                             | Description                                                 | Status      |
| --- | ------------------------------------------------- | ----------------------------------------------------------- | ----------- |
| 01  | [Context Primary Keys](./context-primary-keys.md) | Creating and managing context primary keys for agent memory | ✅ Complete |
| 02  | [Memory Strategies](./memory-strategies.md)       | Lazy, Hydration, and MCP approaches                         | ✅ Complete |
| 03  | [Memory Hydration](./memory-hydration.md)         | Pre-loading knowledge via synthetic conversations           | ✅ Complete |
| 04  | [Threading Phase](./threading-phase.md)           | Context detection and key identification                    | ✅ Complete |

---

## Quick Links

### For Administrators

- [How to create context primary keys](./context-primary-keys.md#creating-context-primary-keys)
- [How to hydrate agent memory](./memory-hydration.md)

### For Developers

- [API Endpoints Reference](../../../developers/api-reference/intro.md#endpoints)
- [Bulk key creation](../../../developers/api-reference/context-primary-keys.md#bulk-provisioning-example-python)

### For Business Users

- [Understanding context primary keys](./context-primary-keys.md#what-is-a-context-primary-key)
- [Designing keys for your domain](./context-primary-keys.md#designing-context-primary-keys)
- [Memory strategies overview](./memory-strategies.md)

---

## Key Concepts at a Glance

### Context Primary Key

A unique identifier that creates an isolated memory space for an agent. Example: contract number `C45653` groups all conversations about that specific contract.

### Memory Strategy

The approach used to populate an agent's memory:

- **Lazy** - Memory grows organically through real conversations
- **Hydration** - Memory is pre-loaded with synthetic past conversations
- **MCP** - Memory is refreshed dynamically from external data sources

### Threading

The process by which an expert agent identifies the relevant context primary key when first mentioned in a conversation.

### Context Window

The LLM's working memory containing the full conversation history for a given context primary key.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     MEMORY                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐   │
│  │Context Primary│  │Context Primary│  │Context Primary│   │
│  │  Key: C45653  │  │  Key: C78901  │  │  Key: C12345  │   │
│  └───────┬───────┘  └───────┬───────┘  └───────┬───────┘   │
│          │                  │                  │            │
│          ▼                  ▼                  ▼            │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐   │
│  │ Conversation  │  │ Conversation  │  │ Conversation  │   │
│  │   History     │  │   History     │  │   History     │   │
│  │               │  │               │  │               │   │
│  │ May 2023      │  │ June 2023     │  │ Sept 2025     │   │
│  │ June 2023     │  │ April 2024    │  │ Oct 2025      │   │
│  │ April 2024    │  │ July 2025     │  │               │   │
│  │ July 2025     │  │ Oct 2025      │  │               │   │
│  │ Oct 2025      │  │               │  │               │   │
│  └───────────────┘  └───────────────┘  └───────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Getting Started

1. **Understand your use case** - Identify what context primary keys make sense for your agent's domain
2. **Create your keys** - Use the UI for small scale, API for bulk operations
3. **Choose a memory strategy** - Lazy for organic growth, Hydration for pre-loaded knowledge
4. **Configure threading** - Set up key detection prompts for expert agents
5. **Test and iterate** - Verify memory works as expected in conversations

---

## Contributing

This documentation is maintained as part of the 6022 Agentic Teams project. For updates or corrections, please submit a pull request to the [documentation repository](https://github.com/6022-labs/documentation).
