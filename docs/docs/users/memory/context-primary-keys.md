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

Need to import thousands of keys or sync them from an external system? Share your requirements with a developer and point them to the [Programmatic Context Primary Keys guide](../../developers/api-reference/context-primary-keys.md).

Provide them with:

- The list of keys to create (CSV or spreadsheet works great)
- The agent collection address + token ID (visible in the orchestrator UI)
- Any validation rules or naming conventions they should enforce

---

## Key Detection: The Threading Phase

When an expert agent is first mentioned in a conversation, it runs a "threading" routine to identify which context primary key applies. The full step-by-step walkthrough (flowchart, key-identification methods, MCP lookup) lives in **[Threading Phase](./threading-phase.md)** — that page is the canonical source.

The short version: the agent runs its **Context Not Found** prompt, attempts to extract or match a context primary key from the message, and if it succeeds the orchestrator persists a row in `agent_context_conversation_mappings` linking the conversation to that key for the rest of its lifetime. There is no mutex or distributed lock — concurrent threading attempts on the same conversation are not guarded.

Regex patterns like the ones in the section above (UUID, contract number, email) are used **inside the prompt template** to teach the agent what a valid key looks like — they are not enforced by the API. `CreateAgentContextRequest.Validate()` only checks that the key is non-empty and does not start with the reserved `conversation-context-discovery/` prefix.

For prompt configuration examples (Context Not Found, Context Found, regex helpers, MCP-assisted lookup), see [Expert Prompts](../prompting/expert-prompts.md).

---

## Managing Existing Context Primary Keys

### Viewing Keys

**Via UI:** Actions → Manage contexts → Browse with pagination

**Via API:**

```http
GET /agent/{agentCollectionAddress}/{agentTokenId}/contexts/keys
```

### Database Direct Access

For advanced troubleshooting, administrators can connect directly to the agent's PostgreSQL database. Each agent has its **own isolated database** — there is no `agent_id` column on `agent_contexts` because the discriminator is the database connection itself.

```sql
-- Connect to the agent's own database, then:
SELECT key FROM agent_contexts;
SELECT COUNT(*) FROM agent_contexts;
```

Internal threading scratchpads stored under the `conversation-context-discovery/` prefix will appear in raw queries but are filtered out of the user-facing API response.

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
