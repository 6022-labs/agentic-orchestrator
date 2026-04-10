---
title: Memory Strategies
sidebar_position: 4
---

# Memory Strategies

## Introduction

Once a context primary key is created, you need to decide **how** the agent's memory will be populated. The 6022 framework offers three distinct strategies, each suited to different use cases.

This document explains when and how to use each strategy.

---

## Strategy Overview

```
                    ┌─────────────────────────┐
                    │  Create Context Primary │
                    │          Key            │
                    └───────────┬─────────────┘
                                │
        ┌───────────────────────┼───────────────────────┐
        │                       │                       │
        ▼                       ▼                       ▼
┌───────────────┐       ┌───────────────┐       ┌───────────────┐
│ LAZY STRATEGY │       │   HYDRATION   │       │ MCP STRATEGY  │
│               │       │   STRATEGY    │       │               │
│  🛋️ Passive   │       │  💉 Proactive │       │  🔄 Dynamic   │
└───────────────┘       └───────────────┘       └───────────────┘
        │                       │                       │
        ▼                       ▼                       ▼
   Conversations           Synthetic              External
   feed memory             "past" chats           data sources
        │                       │                       │
        ▼                       ▼                       ▼
┌───────────────┐       ┌───────────────┐       ┌───────────────┐
│  Happiness    │       │  GDPR Agent   │       │  CRM Agent    │
│   Manager     │       │               │       │               │
└───────────────┘       └───────────────┘       └───────────────┘
```

| Strategy      | Description                                            | Best For                                |
| ------------- | ------------------------------------------------------ | --------------------------------------- |
| **Lazy**      | Memory grows organically through real conversations    | General support, relationship building  |
| **Hydration** | Memory is pre-loaded with synthetic past conversations | Compliance, knowledge bases, onboarding |
| **MCP**       | Memory is refreshed dynamically from external systems  | CRM, ERP, real-time data needs          |

---

## Lazy Strategy

### Concept

The **lazy strategy** is the simplest approach: the agent's memory grows naturally through actual conversations with humans and other agents. No proactive initialization is required.

The agent starts with an empty memory for each context primary key. As conversations occur, the history accumulates organically.

### Characteristics

- ✅ No setup required beyond creating context primary keys
- ✅ Memory reflects genuine interactions
- ✅ Low maintenance overhead
- ⚠️ Agent has no knowledge until first interaction
- ⚠️ May need time to build useful context

### When to Use

| Use Case            | Why Lazy Works                                     |
| ------------------- | -------------------------------------------------- |
| Happiness Manager   | Customer relationships develop naturally over time |
| General Support     | Issues are discussed as they arise                 |
| Team Collaboration  | Context builds through ongoing exchanges           |
| Feedback Collection | Responses accumulate organically                   |

### Implementation

1. Create the context primary key (see [Context Primary Keys](./context-primary-keys.md))
2. That's it — conversations will automatically populate memory

### Example Flow

```
Day 1: Customer C45653 contacts support
       → Agent has no history, starts fresh conversation
       → Conversation is stored under key C45653

Day 15: Customer C45653 returns with follow-up
        → Agent recalls previous conversation
        → Provides contextual response

Day 45: Customer C45653 has new issue
        → Agent has full history of both previous interactions
        → Can reference past solutions and preferences
```

---

## Hydration Strategy

### Concept

The **hydration strategy** pre-populates an agent's memory with synthetic "past" conversations. These fabricated exchanges give the agent knowledge it needs from day one.

The agent genuinely believes it had these conversations. When the LLM processes future requests, the hydrated content appears in its context window as real history.

### Characteristics

- ✅ Agent has immediate expertise
- ✅ Consistent baseline knowledge across all keys
- ✅ Can embed procedures, guidelines, regulations
- ⚠️ Requires upfront preparation
- ⚠️ May need periodic updates as source information changes
- ⚠️ Agent may show confusion if challenged about hydrated memories

### When to Use

| Use Case             | Why Hydration Works                            |
| -------------------- | ---------------------------------------------- |
| GDPR Agent           | Needs complete regulatory knowledge from start |
| Compliance Agent     | Must know all procedures immediately           |
| Onboarding Agent     | Requires company policies pre-loaded           |
| Knowledge Base Agent | Domain expertise needed before first query     |

### Implementation

See [Memory Hydration](./memory-hydration.md) for collaboration tips and share the linked developer guide when you need the API details.

**High-Level Process:**

1. Create the context primary key
2. Design the synthetic conversation(s) to inject
3. Use the API to add messages with appropriate roles
4. Agent now "remembers" this information

### Example Hydration

For a GDPR Agent with context primary key `Article-15`:

```
Synthetic conversation to inject:

User: "Note for later the whole content of GDPR chapter 'individual privacy 3.1':
       The right of access allows data subjects to obtain confirmation as to whether
       personal data concerning them is being processed..."

Agent: "Yes, duly noted. I will reference this information when handling
        Article 15 requests."
```

The agent now has this regulatory knowledge "in memory" and will use it naturally in responses.

### Hydration Side Effects

When an agent's memory is hydrated with fabricated information, it may occasionally show inconsistencies:

```
User: "Are you serious about what you said earlier?"

Agent: "Sorry — I don't know why I said that earlier, I wasn't fully awake..."
```

This is expected behavior when the agent cannot truly recall the context of hydrated memories.

---

## MCP Strategy

### Concept

The **MCP (Model Context Protocol) strategy** connects agents to external data sources that dynamically provide and refresh memory content. Instead of static knowledge, the agent accesses live data from CRMs, databases, APIs, and other systems.

### Characteristics

- ✅ Always up-to-date information
- ✅ Integrates with existing business systems
- ✅ Scales to large, changing datasets
- ⚠️ Requires MCP server setup and maintenance
- ⚠️ Dependent on external system availability
- ⚠️ May have latency for complex queries

### When to Use

| Use Case         | Why MCP Works                                              |
| ---------------- | ---------------------------------------------------------- |
| CRM Agent        | Needs current customer data from HubSpot, Salesforce, etc. |
| Inventory Agent  | Must reflect real-time stock levels                        |
| Scheduling Agent | Requires live calendar availability                        |
| Analytics Agent  | Needs current metrics and KPIs                             |

### How MCP Works

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│    Agent    │────▶│  MCP Server │────▶│  External   │
│             │◀────│             │◀────│   System    │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   CRM       │
                    │   ERP       │
                    │   Database  │
                    │   API       │
                    └─────────────┘
```

1. Agent receives a query requiring external data
2. Agent calls the configured MCP server
3. MCP server fetches data from external systems
4. Data is returned and incorporated into agent's response
5. Data is automatically stored in agent memory for future reference

### MCP Configuration

In the orchestrator UI:

1. Navigate to your agent
2. Open **Actions** → **Manage MCPs**
3. Configure the MCP server connection:
   - Server URL
   - Authentication (Bearer token, API key)
   - Allowed tools/endpoints

### MCP During Threading

MCP servers are accessible during the threading phase, allowing agents to:

- Validate context primary keys against external databases
- Look up customer information from partial identifiers
- Cross-reference data before committing to a key

---

## Combining Strategies

Strategies are not mutually exclusive. Many production deployments combine approaches:

### Hydration + Lazy

**Use case:** Agent needs baseline knowledge but also learns from interactions

```
1. Hydrate with core procedures and guidelines
2. Let organic conversations add case-specific context
```

**Example:** Customer support agent hydrated with return policy, then learns individual customer preferences through conversations.

### Hydration + MCP

**Use case:** Agent needs static knowledge plus dynamic data

```
1. Hydrate with regulatory content, procedures
2. Use MCP for real-time customer/inventory data
```

**Example:** GDPR agent hydrated with regulations, uses MCP to check customer consent status.

### All Three

**Use case:** Complex enterprise deployments

```
1. Hydrate with organizational knowledge
2. Connect MCP for live system data
3. Allow organic learning from team interactions
```

---

## Strategy Selection Guide

```
                              START
                                │
                                ▼
                   ┌────────────────────────┐
                   │ Does agent need        │
                   │ knowledge from day 1?  │
                   └───────────┬────────────┘
                               │
              ┌────────────────┴────────────────┐
              │ NO                              │ YES
              ▼                                 ▼
     ┌─────────────────┐           ┌────────────────────────┐
     │  LAZY STRATEGY  │           │ Is the knowledge       │
     │                 │           │ in external systems?   │
     └─────────────────┘           └───────────┬────────────┘
                                               │
                          ┌────────────────────┴────────────────┐
                          │ NO                                  │ YES
                          ▼                                     ▼
                 ┌─────────────────┐                   ┌─────────────────┐
                 │    HYDRATION    │                   │  MCP STRATEGY   │
                 │    STRATEGY     │                   │                 │
                 └─────────────────┘                   └─────────────────┘
```

---

## Summary

| Strategy      | Setup Effort | Maintenance      | Data Freshness           | Best For                               |
| ------------- | ------------ | ---------------- | ------------------------ | -------------------------------------- |
| **Lazy**      | None         | None             | Real-time (organic)      | Relationship building, general support |
| **Hydration** | Medium       | Periodic updates | Static (until refreshed) | Compliance, procedures, onboarding     |
| **MCP**       | High         | Ongoing          | Real-time (from source)  | CRM, ERP, live data integration        |

Choose based on your agent's mission, data requirements, and operational constraints. Most production deployments benefit from combining strategies.

---

## Related Documentation

- [Context Primary Keys](./context-primary-keys.md) - Creating keys before applying strategies
- [Memory Hydration](./memory-hydration.md) - Detailed API for injecting memories
- [Threading Phase](./threading-phase.md) - How agents detect context
