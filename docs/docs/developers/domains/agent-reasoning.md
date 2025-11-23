---
id: agent-reasoning
title: Domain — Agent Reasoning
sidebar_label: Agent Reasoning
---

This domain is the agent’s thinking stack. It composes prompts from the latest context, talks to LLM providers, triggers MCP tools, and returns responses that channels can send back to humans.

What it owns:

- Model connections (provider choice, model name, API keys, base URLs) so each agent can decide how it wants to reason.
- MCP tool links for everything that needs authenticated JSON-RPC access (search, code, proprietary tools, etc.).
- Chat completion endpoints that other domains call when they want “an answer from this agent.”
- Setup/status state so we know whether an agent is actually ready to think before routing traffic to it.

How it works (high level):

- Receives a prompt + conversation slice from `agent-knowledge`, picks the configured model for that agent, and forwards the request to the provider over HTTP.
- Normalizes provider responses, tracks whether the reasoning round is final, and hands the text back to the caller (usually `channels-and-integrations`).
- Keeps secrets (LLM API keys, MCP headers) in secure storage so reasoning calls can be executed without re-entering credentials.
- Pushes responses and conversation in general during each reasoning round to `agent-knowledge` to reuse at next turn.
