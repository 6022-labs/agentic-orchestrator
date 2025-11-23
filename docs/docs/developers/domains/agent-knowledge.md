---
id: agent-knowledge
title: Domain — Agent Knowledge
sidebar_label: Agent Knowledge
---

This domain keeps each agent’s long-term memory. It provisions private databases, stores prompt templates and conversation contexts, and serves them back to the reasoning stack before every inference. Nothing is on-chain; all state lives in the agent’s own database.

What it owns:

- Per-agent data stores so knowledge never leaks between tenants.
- Prompt templates, message history, and conversation-to-context links that describe “what the agent knows right now.”
- Connection metadata (where the database lives, how to reach it, which migrations have run).
