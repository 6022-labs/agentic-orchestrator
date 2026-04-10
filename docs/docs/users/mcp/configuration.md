---
sidebar_position: 2
---

# MCP Configuration

This guide explains how to add and configure MCP connections for your agents using the orchestrator interface.

## Prerequisites

Before adding an MCP connection:

1. ✅ Your agent must use a compatible LLM (OpenAI, Anthropic, or AWS Bedrock)
2. ✅ You have the MCP server URL (must be HTTPS)
3. ✅ You have authentication credentials (if required)
4. ✅ The MCP server is running and accessible

## Adding an MCP Connection

### Step 1: Access MCP Settings

1. Go to **Dashboard** → **Agents**
2. Find your agent
3. Click the **⋮** menu
4. Select **"Manage MCPs"**

### Step 2: Add New Connection

1. Click **"+ Add MCP Connection"**
2. Fill in the configuration form

### Step 3: Configure Connection Details

| Field | Description | Required |
|-------|-------------|:--------:|
| **Name** | Identifier for this MCP (e.g., `stripe-mcp`). Must match `^[a-z][a-z0-9-]*$` and be at least 3 characters. Must be unique per agent. | ✅ |
| **URL** | Server endpoint (must be HTTPS). Must be unique per agent. | ✅ |
| **Authentication** | Set of HTTP header key/value pairs that get attached to every request. The UI presents this as "None / API Key (Bearer) / Custom Header / Basic Auth" but they all resolve to the same underlying header map. | ✅ |

## URL Requirements

MCP URLs must use **HTTPS** — `http://` URLs are rejected by the orchestrator's validation. URLs are also trimmed of leading/trailing whitespace before validation. For trailing-slash gotchas and hidden-character issues, see [Troubleshooting → URL Issues](./troubleshooting.md#url-issues).

## Authentication

Authentication is a set of HTTP header key/value pairs. Whatever the UI labels them, they all resolve to headers that the orchestrator attaches to every MCP request. Here is what each label produces:

| UI label | What you fill in | What the orchestrator sends |
|---|---|---|
| **None** | nothing | no auth header |
| **API Key (Bearer)** | a token | `Authorization: Bearer <token>` |
| **Custom Header** | a header name + value | `<header-name>: <value>` |
| **Basic Auth** | username + password | `Authorization: Basic <base64(user:pass)>` |

Header values are stored in the orchestrator's secret vault, not in the agent database. When you delete an MCP connection, the vault entries for its header values are removed alongside the database row.

## Configuration Examples

For pre-configured server examples (Stripe, Kirha, Markitdown, HubSpot, Notion, etc.), see **[Available MCPs](./available-mcps.md)** — that page is the canonical catalogue. Don't duplicate the per-server examples here.

## Managing Multiple MCPs

An agent can have multiple MCP connections. Each connection:

- Has a unique name
- Has its own URL and authentication
- Provides its own set of tools

The LLM sees tools from ALL connected MCPs and decides which to use.

### Recommended Approach

Instead of one large MCP with everything, use multiple specialized MCPs:

```
Agent: customer-service.mycompany.137.protocol6022.eth
├── MCP: crm        → Customer data, history
├── MCP: payments   → Stripe invoices, subscriptions
├── MCP: documents  → Contract templates, policies
└── MCP: websearch  → Real-time information
```

## Verifying Your Configuration

### Test with Postman or cURL

Before adding to the orchestrator, verify the MCP server responds correctly:

```bash
curl -X POST 'https://your-mcp-server.com/mcp' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_API_KEY' \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "params": {},
    "id": 1
  }'
```

**Expected response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "tools": [
      {
        "name": "tool_name",
        "description": "What this tool does",
        "inputSchema": { ... }
      }
    ]
  },
  "id": 1
}
```

### Using MCP Inspector

For a visual tool discovery experience, use [MCP Inspector](https://github.com/modelcontextprotocol/inspector):

```bash
npx @modelcontextprotocol/inspector https://your-mcp-server.com/mcp
```

## Removing an MCP Connection

1. Go to **Manage MCPs**
2. Find the connection to remove
3. Click the **⋮** menu on that row
4. Select **"Delete"**
5. Confirm removal

:::note
Removing an MCP does not affect past conversations. The agent simply won't have access to those tools in future interactions.
:::

## Best Practices

For specialization strategies, tool count limits, and the full Do's/Don'ts list (when to split MCPs, when to combine, how to prompt for tool usage), see **[MCP Best Practices](./best-practice.md)**.