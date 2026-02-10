---
sidebar_position: 3
---

# Available MCP Servers

This page lists MCP servers that have been tested with 6022 Agentic Teams, along with their configuration details and known issues.

## Production-Ready MCPs

These MCPs have been tested in production environments and work reliably.

### Stripe MCP

**Purpose:** Access Stripe payment data (customers, invoices, subscriptions, balances)

| Field | Value                                                   |
| ----- | ------------------------------------------------------- |
| URL   | `https://mcp.stripe.com`                                |
| Auth  | API Key (Bearer)                                        |
| Key   | Your Stripe Secret Key (`sk_live_...` or `sk_test_...`) |

**Available Tools:**

- List customers
- Get customer details
- List invoices
- Get subscription status
- Query payment intents

**Example Use:**

> "What are the last 5 invoices for customer john@example.com?"

The agent will call Stripe MCP to fetch invoice data.

---

### Kirha MCP

**Purpose:** Web search, crypto data, AI-powered research

| Field        | Value                                       |
| ------------ | ------------------------------------------- |
| URL          | `https://kirha-mcp-6022.up.railway.app/mcp` |
| Auth         | Custom Header                               |
| Header Key   | `x-kirha-api-key`                           |
| Header Value | Your Kirha API Key                          |

**Available Tools:**

- `websearch` — Search the internet for current information
- `crypto_prices` — Get real-time cryptocurrency prices
- `crypto_analysis` — Market analysis and indicators

**Getting a Kirha API Key:**
Contact Kirha or your 6022 administrator to obtain API credits.

:::note Credits Required
Kirha tools consume credits. Check your balance if you get 500 errors.
:::

---

### Markitdown MCP

**Purpose:** Convert documents (PDF, DOCX, XLSX, images) to Markdown text

| Field | Value                                                   |
| ----- | ------------------------------------------------------- |
| URL   | `https://markitdown-mcp-production.up.railway.app/mcp/` |
| Auth  | None                                                    |

:::warning Trailing Slash Required
The URL **must** end with `/mcp/` (with slash). Using `/mcp` will fail.
:::

**Available Tools:**

- `convert_url` — Convert a document from URL to Markdown
- `convert_base64` — Convert a base64-encoded document to Markdown

**Supported Formats:**

- PDF (structured documents)
- DOCX (Word documents)
- XLSX (Excel spreadsheets)
- Images (extracts metadata)

**Limitations:**

- Scanned PDFs may return limited data
- Very large files may timeout
- Cannot follow redirects well

**Example Use:**

> "Extract the content from this PDF: https://example.com/contract.pdf"

---

## Beta / Limited MCPs

These MCPs work but have known limitations or require special setup.

### HubSpot Remote MCP

**Purpose:** Access HubSpot CRM data

| Field | Value                                 |
| ----- | ------------------------------------- |
| URL   | `https://app.hubspot.com/mcp/v1/http` |
| Auth  | API Key (Bearer)                      |
| Key   | HubSpot Private App Token             |

:::caution Private Beta
As of late 2025, HubSpot's remote MCP is in private beta. You may need to request access via their form.
:::

**Known Issues:**

- `Invalid tool schema` errors for some tools
- May require specific HubSpot subscription tier
- Rate limiting can be aggressive

**Alternative:** Build your own HubSpot MCP using their REST API.

---

### Notion MCP

**Purpose:** Access Notion pages and databases

| Official URL | `https://mcp.notion.com/mcp`  |
| ------------ | ----------------------------- |
| Auth         | OAuth 2.1 (not simple Bearer) |

:::danger OAuth Limitation
The official Notion MCP uses OAuth 2.1, which is **not compatible** with the 6022 orchestrator's current authentication methods (Bearer token, Custom Header).
:::

**Solution: Self-Host Notion MCP**

Deploy the official Notion MCP server yourself:

```bash
# Clone the official repo
git clone https://github.com/makenotion/notion-mcp-server

# Deploy to Railway, Fly.io, or your server
# Configure with your Notion integration token
```

Then configure in orchestrator:

| Field | Value                                     |
| ----- | ----------------------------------------- |
| URL   | `https://your-notion-mcp.railway.app/mcp` |
| Auth  | API Key (Bearer)                          |
| Key   | Your custom auth token                    |

**Getting a Notion Integration Token:**

1. Go to https://www.notion.so/my-integrations
2. Create a new integration
3. Copy the "Internal Integration Secret" (starts with `ntn_`)
4. Share your Notion pages with the integration

---

## Custom MCPs

You can build and deploy your own MCP servers for company-specific data.

### Example: Company Knowledge Base MCP

Build an MCP that exposes your internal documentation:

```python
# Example tools your MCP could provide
tools = [
    {
        "name": "search_policies",
        "description": "Search company HR and legal policies",
        "inputSchema": {
            "type": "object",
            "properties": {
                "query": {"type": "string"}
            }
        }
    },
    {
        "name": "get_procedure",
        "description": "Get step-by-step procedure for a process",
        "inputSchema": {
            "type": "object",
            "properties": {
                "procedure_name": {"type": "string"}
            }
        }
    }
]
```

### Example: Client-Specific MCP (Insurance)

A Heymondo-style MCP for insurance claims:

| Field    | Value                                      |
| -------- | ------------------------------------------ |
| URL      | `https://hms-mcp.company.com/mcp/heymondo` |
| Auth     | Basic Auth                                 |
| Username | `company_mcp`                              |
| Password | `SecureP@ss2026`                           |

**Tools provided:**

- `get_claim_status` — Check claim processing status
- `get_policy_details` — Retrieve policy information
- `search_documents` — Find related documents

---

## MCP Comparison Table

| MCP                  | Auth Type     | LLM Support                | Status          |
| -------------------- | ------------- | -------------------------- | --------------- |
| Stripe               | Bearer        | OpenAI, Anthropic, Bedrock | ✅ Production   |
| Kirha                | Custom Header | OpenAI, Anthropic, Bedrock | ✅ Production   |
| Markitdown           | None          | All compatible             | ✅ Production   |
| HubSpot              | Bearer        | OpenAI                     | ⚠️ Beta         |
| Notion (official)    | OAuth         | N/A                        | ❌ Incompatible |
| Notion (self-hosted) | Bearer        | OpenAI, Anthropic, Bedrock | ✅ Works        |

---

## Requesting New MCPs

If you need an MCP for a service not listed here:

1. **Check if the service offers an MCP** — Many SaaS platforms are adding MCP support
2. **Build a custom MCP** — See [Building Custom MCPs](./buildings)
3. **Contact 6022 support** — We may be able to help build or recommend solutions

## MCP Resources

- [MCP Specification](https://modelcontextprotocol.io/specification/2025-06-18/server/tools)
- [MCP Inspector](https://github.com/modelcontextprotocol/inspector)
- [Anthropic MCP Documentation](https://docs.anthropic.com/en/docs/mcp)
