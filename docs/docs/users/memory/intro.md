---
title: Introduction
sidebar_position: 1
---

# Memory

Everything an agent "knows" and "remembers" is organised around **context primary keys**. A key groups the conversation history for one specific thing — a customer, a contract, a ticket, a GDPR article — so that when an expert is asked about that thing again, it can load everything it ever discussed under that key.

## Three concepts to understand

- **[Context primary key](./context-primary-keys.md)** — a unique identifier that creates an isolated memory space. Example: contract number `C45653` groups every conversation the CRM expert has ever had about that customer.
- **[Memory strategy](./memory-strategies.md)** — how the memory gets populated. Lazy (grows from real conversations), Hydration (pre-loaded synthetic conversations), or MCP (fetched live from an external system).
- **[Threading phase](./threading-phase.md)** — when an expert is first @-mentioned in a conversation, it runs a threading routine to identify which context primary key applies and lock it for the rest of the thread.

## Getting started

1. Decide what makes a good key for your agent's domain (email? contract number? ticket ID?)
2. Create the keys through the dashboard for small batches, or through the API for bulk
3. Pick a memory strategy that matches your use case
4. Configure the agent's Context Not Found prompt to recognise your key format
5. Test in a real conversation

## For developers

- [API Reference Overview](../../developers/api-reference/intro.md)
- [Context Primary Keys API](../../developers/api-reference/context-primary-keys.md) — bulk provisioning, pagination, search
- [Memory Hydration API](../../developers/api-reference/memory-hydration.md)
