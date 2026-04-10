---
sidebar_position: 1
---

# Introduction to MCP

MCP (Model Context Protocol) is an open standard that lets agents call external tools — CRMs, databases, file converters, web search, payment APIs, or anything else with an HTTP endpoint. Each agent can connect to multiple MCP servers simultaneously; the LLM sees all the available tools and picks the right one for each question.

```
┌─────────┐   JSON-RPC 2.0   ┌─────────────┐   calls   ┌──────────────┐
│  Agent  │ ───────────────▶ │ MCP Server  │ ───────── │   External   │
│  (LLM)  │  tools/list      │             │           │    System    │
│         │  tools/call      │             │           │              │
└─────────┘ ◄─────────────── └─────────────┘ ◄──────── └──────────────┘
```

## LLM compatibility

The orchestrator distinguishes **native MCP routing** (server definitions are forwarded directly to the provider) from **function-call loop routing** (the orchestrator runs its own loop and asks the model to emit tool calls). The canonical per-provider matrix lives in [LLM Introduction → Supported Providers](../llm/intro.md#supported-providers); the short version:

| Provider | MCP routing |
|---|---|
| **OpenAI** | ✅ Native (via `/v1/responses`) |
| **Anthropic** | ✅ Native (requires custom-header auth) |
| **AWS Bedrock** | ⚙️ Function-call loop (Converse API) |
| **xAI / Mistral / Deepseek / LibertAI / Ollama** | ⚙️ Function-call loop — depends on the model's function-calling support |

If you configure MCP on an agent whose model doesn't actually support function calling, the orchestrator still attempts to invoke the loop; tools will simply never be selected.

## What MCPs look like in practice

MCP uses JSON-RPC 2.0 over HTTPS with two methods: `tools/list` (discover what a server offers) and `tools/call` (execute one). You can see the exact wire format by testing a server with cURL or MCP Inspector — see [Troubleshooting](./troubleshooting.md) for the exact commands.

Common use cases:

- **Live data** — CRM lookups (HubSpot), payments (Stripe), scheduling, inventory
- **Document processing** — PDF/DOCX → Markdown (Markitdown)
- **Research** — web search (Kirha), crypto market data, knowledge bases

## Authentication model

Each MCP connection stores a set of HTTP header key/value pairs that get attached to every request. The UI labels these as "None / API Key (Bearer) / Custom Header / Basic Auth" for clarity, but under the hood they all resolve to the same header map. The values are stored in the orchestrator's secret vault, never in plaintext in the database.

## Next steps

- **[Configuration](./configuration.md)** — set up a connection in the dashboard
- **[Available MCPs](./available-mcps.md)** — pre-built servers (Stripe, Kirha, Markitdown, etc.)
- **[Best practices](./best-practice.md)** — how to structure MCPs for a swarm
- **[Troubleshooting](./troubleshooting.md)** — HTTP errors, URL gotchas, testing

**Building your own MCP server?** That's a developer task — see [Building a Custom MCP Server](../../developers/integrations/building-mcps.md) for the Python/Flask template and deployment guide.
