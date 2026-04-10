---
title: Memory Strategies
sidebar_position: 4
---

# Memory Strategies

Once a context primary key exists, you have to decide **how** its memory gets populated. Three strategies, and you can combine them.

| Strategy | How it works | Best for | Setup effort |
|---|---|---|---|
| **Lazy** | Memory grows organically from real conversations | Relationship building, general support | None |
| **Hydration** | Memory is pre-loaded with synthetic "past" conversations via API | Compliance, procedures, onboarding knowledge | Medium |
| **MCP** | Memory is fetched live from an external system on every query | CRM, ERP, real-time data needs | High |

## Lazy

The simplest approach: do nothing. The agent starts with empty memory for each key, and builds up history naturally through real conversations.

- **Day 1**: customer `C45653` contacts support. Agent has no history. Conversation is stored.
- **Day 15**: customer returns. Agent recalls day-1 conversation and responds with context.
- **Day 45**: third issue. Agent has the full history of both previous interactions.

Use this for happiness managers, general support, team collaboration — anywhere relationships develop over time.

**Pros:** zero setup; memory reflects actual conversations. **Cons:** the agent is cold on day 1.

## Hydration

Pre-populate memory with synthetic conversations before the agent ever talks to a real user. The agent treats the hydrated content as real history — the LLM sees it in its context window and acts on it.

Example for a GDPR agent with key `Article-15`:

```
User: "Note for later the content of GDPR chapter individual privacy 3.1:
       The right of access allows data subjects to obtain confirmation..."

Agent: "Noted. I will reference this when handling Article 15 requests."
```

Now when a real question arrives for `Article-15`, the agent "remembers" the regulation. Use this for compliance, onboarding, knowledge bases — anywhere the agent needs to know things from day one.

**Pros:** immediate expertise; consistent baseline across all keys. **Cons:** requires upfront work; the agent may show minor inconsistencies if directly challenged about fabricated memories.

See [Memory Hydration](./memory-hydration.md) for the conceptual patterns and the [Memory Hydration API](../../developers/api-reference/memory-hydration.md) for the runnable code.

## MCP

Instead of storing knowledge in memory, give the agent a tool that fetches it live from an external system. The agent calls the MCP server on every question that needs fresh data.

Use this when the data changes faster than you can update it — real-time customer data from HubSpot, live inventory, current appointment slots, current KPIs. The agent is always up-to-date but depends on the external system being available.

**Pros:** always fresh; scales to huge datasets; integrates with existing systems. **Cons:** requires MCP server setup and maintenance; latency per query.

See [MCP Configuration](../mcp/configuration.md) and [MCP Best Practices](../mcp/best-practice.md) for how to wire one up.

## Combining strategies

Most production deployments combine two or three:

- **Hydration + Lazy** — hydrate baseline procedures, then let organic conversations add case-specific context. Good for customer support agents with a return policy + individual customer preferences.
- **Hydration + MCP** — hydrate static regulatory knowledge, use MCP for dynamic customer data. Good for a GDPR agent that also needs to check current consent status.
- **All three** — hydrate internal docs, MCP for live system data, lazy for organic team interactions. Complex enterprise deployments.

## Picking a strategy

Ask two questions:

1. **Does the agent need knowledge on day 1?** If no → **Lazy**. If yes, continue.
2. **Is that knowledge in an external system that updates in real time?** If yes → **MCP**. If no → **Hydration**.
