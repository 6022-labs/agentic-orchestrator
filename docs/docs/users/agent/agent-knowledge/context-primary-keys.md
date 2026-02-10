---
title: Context Primary Keys
sidebar_position: 2
---

# Context Primary Keys

## Introduction

In the 6022 Agentic Teams framework, a **context primary key** is the fundamental mechanism that organizes an agent's memory. Each context primary key creates an isolated memory space where all related conversations are stored and retrieved.

This document covers the creation and management of context primary keys for expert agents.

---

## What is a Context Primary Key?

A context primary key is a unique identifier that:

- Links an agent to a specific conversation context
- Groups all related exchanges under one retrievable memory
- Enables the agent to recall past interactions when the same key is referenced

**Example:** A CRM agent with context primary key `C45653` will store and retrieve all conversations about customer C45653 separately from conversations about customer `C78901`.

---

## Context Primary Key Behavior by Agent Type

| Agent Type       | Context Primary Key                  | Creation                                       |
| ---------------- | ------------------------------------ | ---------------------------------------------- |
| **Facilitator**  | The conversation itself              | Automatic - each new conversation = new memory |
| **Human**        | Set of ETH addresses of participants | Automatic - known upon channel entry           |
| **Digital Twin** | Set of ETH addresses of participants | Automatic - inherited from cloned human        |
| **Expert**       | Business identifier                  | **Manual** - must be explicitly created        |

**This document focuses on Expert agents**, which require manual context primary key management.

---

## Designing Context Primary Keys

### Align with Business Logic

The context primary key must directly relate to the agent's mission. Choose identifiers that naturally map to your domain:

| Agent Mission         | Recommended Context Primary Key | Example Values                       |
| --------------------- | ------------------------------- | ------------------------------------ |
| Customer Relationship | Customer email or ID            | `customer@example.com`, `CUST-45653` |
| Claims Management     | Claim number                    | `CLM-2025-00142`                     |
| GDPR Compliance       | GDPR chapter/article            | `Article-15`, `Chapter-3.1`          |
| Accounting            | Invoice number                  | `INV-2025-0001`                      |
| Customer Support      | Contract number                 | `312123`, `C45653`                   |
| Technical Support     | Ticket ID                       | `TKT-8889f6a1`                       |

### Key Design Rules

**1. Maintain Consistency**

Never mix different key types within the same agent. A CRM agent should not have both customer IDs and GDPR chapters as context primary keys.

```
✅ Good: CRM Agent keys = [CUST-001, CUST-002, CUST-003]
❌ Bad:  CRM Agent keys = [CUST-001, Article-15, INV-2025]
```

**2. Use Predictable Formats**

Define a clear pattern for your context primary keys to enable validation:

| Format Type     | Regex Pattern                                                               | Example                                |
| --------------- | --------------------------------------------------------------------------- | -------------------------------------- |
| UUID            | `^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$` | `8889f6a1-df3e-423f-a99c-2e5ce120c489` |
| Contract Number | `^C[0-9]{5}$`                                                               | `C45653`                               |
| Email           | Standard email regex                                                        | `contact@example.com`                  |
| Custom Prefix   | `^[A-Z]{3}-[0-9]{4}-[0-9]{4}$`                                              | `CLM-2025-0142`                        |

**3. Plan for Scale**

Consider your expected volume:

- Small scale (< 100 keys): Manual creation via UI is practical
- Medium scale (100 - 50,000 keys): Use API with scripts

---

## Creating Context Primary Keys

### Method 1: Administration Interface

For small-scale or ad-hoc key creation:

1. Navigate to your agent in the orchestrator UI
2. Open the **Actions** menu (⋮)
3. Select **Manage contexts**
4. Click **+ Add Agent Context**
5. Enter the context primary key value
6. Click **Save**

The interface also allows you to:

- View all existing context primary keys with pagination
- Search/filter keys (planned feature)

### Method 2: Automations (Ask Engineering)

Need to import thousands of keys or sync them from an external system? Share your requirements with a developer and point them to the [Programmatic Context Primary Keys guide](../../../developers/api-reference/context-primary-keys.md).

Provide them with:

- The list of keys to create (CSV or spreadsheet works great)
- The agent collection address + token ID (visible in the orchestrator UI)
- Any validation rules or naming conventions they should enforce

---

## Key Detection: The Threading Phase

When an expert agent is mentioned in a conversation, it must identify which context primary key applies. This is called the **threading phase**.

### How Threading Works

1. Agent receives a message and enters `THREADING` state
2. The **Context Not Found** prompt analyzes the message
3. Agent attempts to extract/match a context primary key
4. If found → key is locked for the conversation
5. If not found → agent can ask other agents for help

### Configuring Key Detection

In the agent's **Manage prompts** settings, configure the **Context Not Found** prompt:

**Example for Contract-Based Keys:**

```
Analyze this message: "{{.LatestRequest.Content}}"

Identify the contract number using this regex pattern: ^C[0-9]{5}$

If you find a valid contract number, use it as the context primary key.
If you cannot identify a contract number, ask the facilitator for clarification.
```

**Example for UUID-Based Keys:**

```
The context primary key follows this UUID format:
^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$

Extract the UUID from the message if present.
If not found, query the MCP to verify the customer exists.
```

### Using MCP for Key Validation

MCP servers are accessible during threading, enabling:

- Validation of context primary keys against external databases
- Lookup of customer IDs from partial information
- Cross-referencing with CRM or ERP systems

---

## Managing Existing Context Primary Keys

### Viewing Keys

**Via UI:** Actions → Manage contexts → Browse with pagination

**Via API:**

```http
GET /agent/{agentCollectionAddress}/{agentTokenId}/contexts/keys
```

### Database Direct Access

For advanced operations, administrators can access the PostgreSQL database:

```sql
-- View all contexts for an agent
SELECT * FROM agent_contexts WHERE agent_id = 'agent-uuid';

-- Count context primary keys
SELECT COUNT(*) FROM agent_contexts WHERE agent_id = 'agent-uuid';
```

**Note:** Direct database access is strongly discouraged and should only be used in exceptional cases for troubleshooting.

---

## Troubleshooting

> Most API-related errors require help from an engineer. Use this table to triage, then share the details with the person running your automations.

| Error                                 | Cause                                  | Solution                                 |
| ------------------------------------- | -------------------------------------- | ---------------------------------------- |
| `context with key X already exists`   | Key was previously created             | Skip or use existing key                 |
| `You are not the owner of this agent` | JWT signed with wrong wallet           | Regenerate JWT with owner's private key  |
| `500 Internal Server Error`           | Database or sync issue                 | Check orchestrator logs, contact support |
| Threading keeps failing               | Key pattern doesn't match prompt regex | Update Context Not Found prompt          |
| Keys not appearing in UI              | Sync delay                             | Wait a few seconds, refresh page         |

---

## Best Practices Summary

1. **Design keys before deployment** - Define your context primary key format and naming convention upfront
2. **Use consistent patterns** - Stick to one key type per agent
3. **Validate key format** - Include regex validation in threading prompts
4. **Plan for bulk operations** - Loop in engineering for API-powered imports when you exceed ~100 keys
5. **Reset threading for demos** - Clear the THREADING context before presentations
6. **Document your key schema** - Keep a reference of what each context primary key format means

---

## Related Documentation

- [Agent Memory Strategies](./memory-strategies.md) - How to populate memory once keys exist
- [Memory Hydration](./memory-hydration.md) - Pre-loading agent knowledge
