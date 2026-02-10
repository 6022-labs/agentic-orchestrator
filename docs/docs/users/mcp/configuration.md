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
| **Name** | Identifier for this MCP (e.g., "stripe-mcp") | ✅ |
| **URL** | Server endpoint (must be HTTPS) | ✅ |
| **Authentication** | Auth method: None, API Key, or Custom Header | ✅ |

## URL Requirements

:::warning Critical
MCP URLs must use **HTTPS**. HTTP is not accepted for security reasons.
:::

**URL Format Rules:**

| ✅ Valid | ❌ Invalid |
|----------|-----------|
| `https://mcp.stripe.com` | `http://mcp.stripe.com` |
| `https://kirha-mcp.up.railway.app/mcp` | `kirha-mcp.up.railway.app/mcp` |
| `https://example.com/mcp/` | `https://example.com/mcp\t` (has tab) |

:::tip Trailing Slash
Some MCP servers require a trailing slash (e.g., `/mcp/` instead of `/mcp`). If you get errors, try adding or removing the trailing slash.
:::

## Authentication Methods

### No Authentication

Some MCP servers are public and don't require authentication:

```
URL: https://markitdown-mcp-production.up.railway.app/mcp/
Authentication: None
```

### API Key (Bearer Token)

The most common method. Creates an `Authorization: Bearer <token>` header automatically:

```
URL: https://mcp.stripe.com
Authentication: API Key
API Key: sk_live_xxxxxxxxxxxxxxxxxxxxx
```

Under the hood, this sends:
```
Authorization: Bearer sk_live_xxxxxxxxxxxxxxxxxxxxx
```

### Custom Header

For MCP servers that require a different header name:

```
URL: https://kirha-mcp-6022.up.railway.app/mcp
Authentication: Custom Header
Header Key: x-kirha-api-key
Header Value: 3YL2MGbf2Gtf6M3GpH1xf_qu6vVfQoyR-Nw7MqDIrDw
```

### Basic Authentication

For servers using username/password (encoded in Base64):

```
URL: https://hms-mcp.example.com/mcp/heymondo
Authentication: Basic Auth
Username: heymondo_mcp
Password: MyP@ssw0rd2026
```

## Configuration Examples

### Stripe MCP

Connect to Stripe for payment data:

| Field | Value |
|-------|-------|
| Name | `stripe` |
| URL | `https://mcp.stripe.com` |
| Authentication | API Key |
| API Key | `sk_live_xxxxx...` (your Stripe secret key) |

### Kirha MCP

Connect to Kirha for web search and crypto data:

| Field | Value |
|-------|-------|
| Name | `kirha` |
| URL | `https://kirha-mcp-6022.up.railway.app/mcp` |
| Authentication | Custom Header |
| Header Key | `x-kirha-api-key` |
| Header Value | `your-kirha-api-key` |

### Markitdown MCP

Convert documents to Markdown (no auth required):

| Field | Value |
|-------|-------|
| Name | `markitdown` |
| URL | `https://markitdown-mcp-production.up.railway.app/mcp/` |
| Authentication | None |

:::warning Trailing Slash Required
The Markitdown URL **must** end with `/mcp/` (with trailing slash). Using `/mcp` will cause errors.
:::

### Custom HubSpot MCP

If you've built your own HubSpot MCP server:

| Field | Value |
|-------|-------|
| Name | `hubspot` |
| URL | `https://mcp-hub.yourdomain.com/mcp` |
| Authentication | API Key |
| API Key | `Bearer eyJhbGci...` (your JWT token) |

## Managing Multiple MCPs

An agent can have multiple MCP connections. Each connection:

- Has a unique name
- Has its own URL and authentication
- Provides its own set of tools

The LLM sees tools from ALL connected MCPs and decides which to use.

### Recommended Approach

Instead of one large MCP with everything, use multiple specialized MCPs:

```
Agent: customer-service.agent.6022.eth
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

### Do's ✅

- Test MCP connections with Postman before configuring
- Use descriptive names for MCP connections
- Keep credentials secure and rotate them regularly
- Use specialized MCPs instead of one mega-MCP

### Don'ts ❌

- Don't use HTTP URLs (HTTPS required)
- Don't include spaces or tabs in URLs
- Don't share API keys between agents if not necessary
- Don't connect too many MCPs (LLM can get confused with 50+ tools)