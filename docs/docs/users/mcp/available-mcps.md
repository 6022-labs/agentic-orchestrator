---
sidebar_position: 3
---

# Available MCP Servers

Catalogue of MCP servers tested with 6022. For how to plug one in, see [Configuration](./configuration.md).

## Production-ready

### Stripe

Payment data (customers, invoices, subscriptions, payment intents).

| Field | Value |
|---|---|
| URL | `https://mcp.stripe.com` |
| Auth | API Key (Bearer) |
| Key | Your Stripe secret key (`sk_live_...` or `sk_test_...`) |

Example prompt: *"What are the last 5 invoices for customer john@example.com?"*

### Kirha

Web search, crypto market data, AI research.

| Field | Value |
|---|---|
| URL | `https://kirha-mcp-6022.up.railway.app/mcp` |
| Auth | Custom Header |
| Header name | `x-kirha-api-key` |
| Header value | Your Kirha API key |

Tools: `websearch`, `crypto_prices`, `crypto_analysis`. Kirha consumes credits — 500 errors usually mean an empty balance. Contact Kirha or your 6022 administrator for credentials.

### Markitdown

Convert documents (PDF, DOCX, XLSX, images) to Markdown text.

| Field | Value |
|---|---|
| URL | `https://markitdown-mcp-production.up.railway.app/mcp/` |
| Auth | None |

:::warning Trailing slash required
The URL **must** end with `/mcp/`. Without the trailing slash it returns 404.
:::

Tools: `convert_url`, `convert_base64`. Scanned PDFs return limited data; very large files may timeout.

## Beta / limited

### HubSpot

| Field | Value |
|---|---|
| URL | `https://app.hubspot.com/mcp/v1/http` |
| Auth | API Key (Bearer) |
| Key | HubSpot Private App Token |

In HubSpot's private beta as of late 2025 — you may need to request access. Known to return `Invalid tool schema` errors on some tools. Rate limiting is aggressive. If it doesn't fit your needs, build a custom HubSpot MCP against their REST API instead.

### Notion

The **official** Notion MCP uses OAuth 2.1, which the 6022 orchestrator doesn't currently support (only static Bearer/Custom Header auth). You have two options:

- **Self-host** the open-source `notion-mcp-server` and protect it with your own Bearer token. Then add it to the orchestrator as a normal MCP connection. The Notion integration token comes from [notion.so/my-integrations](https://notion.so/my-integrations) — use the "Internal Integration Secret" (`ntn_...`) and share the pages you want accessible with the integration.
- **Wait** for the orchestrator to add OAuth support.

## Comparison

| MCP | Auth | Works with | Status |
|---|---|---|---|
| Stripe | Bearer | OpenAI, Anthropic, Bedrock | ✅ Production |
| Kirha | Custom header | OpenAI, Anthropic, Bedrock | ✅ Production |
| Markitdown | None | All providers | ✅ Production |
| HubSpot | Bearer | OpenAI | ⚠️ Beta, schema issues |
| Notion (official) | OAuth | — | ❌ Not compatible |
| Notion (self-hosted) | Bearer | OpenAI, Anthropic, Bedrock | ✅ Works |

## Need something not on this list?

- Check if the service ships an MCP of its own (many SaaS vendors are adding them).
- **Build your own MCP** — see [Building a Custom MCP Server](../../developers/integrations/building-mcps.md) in the developer docs for the Python/Flask template.
- Contact 6022 support — we may already have something in the works.

## References

- [MCP Specification](https://modelcontextprotocol.io/specification/2025-06-18/server/tools)
- [MCP Inspector](https://github.com/modelcontextprotocol/inspector) — visual tool to debug an MCP server
- [Anthropic MCP documentation](https://docs.anthropic.com/en/docs/mcp)
