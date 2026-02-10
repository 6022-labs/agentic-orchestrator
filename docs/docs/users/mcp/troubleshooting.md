---
sidebar_position: 4
---

# MCP Troubleshooting

This guide covers common MCP issues encountered in production, with real-world solutions based on actual deployments.

## Quick Diagnosis

```
MCP not working?
       │
       ▼
┌──────────────────┐
│ What error code? │
└────────┬─────────┘
         │
    ┌────┴────┬────────┬──────────┐
    ▼         ▼        ▼          ▼
  401       424      500     No response
   │         │        │          │
   ▼         ▼        ▼          ▼
Bad API   tools/list  Server   URL/Network
  Key     failed     error      issue
```

---

## Error 401: Unauthorized

**Error message:**
```json
{
  "error": {
    "message": "Error retrieving tool list from MCP server: 'mcp-name'. Http status code: 401 (Unauthorized)",
    "type": "external_connector_error",
    "param": "tools",
    "code": "http_error"
  }
}
```

**Causes:**
1. API key is invalid
2. API key has expired
3. API key doesn't have required permissions
4. Wrong authentication method selected

**Solutions:**

| Check | Action |
|-------|--------|
| Key validity | Generate a new API key from the service |
| Key format | Ensure no extra spaces before/after the key |
| Auth method | Verify you're using correct auth type (Bearer vs Custom Header) |
| Permissions | Check the key has required scopes/permissions |

**Real example:** Notion MCP returning 401 because Bearer token auth is used but Notion requires OAuth. Solution: Self-host Notion MCP instead.

---

## Error 424: Failed Dependency

**Error message:**
```json
{
  "error": {
    "message": "Error retrieving tool list from MCP server: 'mcplpl'. Http status code: 424 (Failed Dependency)",
    "type": "external_connector_error",
    "param": "tools",
    "code": "http_error"
  }
}
```

**Causes:**
1. MCP server cannot connect to its backend
2. MCP server configuration error
3. URL contains invalid characters

**Solutions:**

| Check | Action |
|-------|--------|
| URL cleanliness | Look for tabs, spaces, or "Open menu" text accidentally copied |
| Backend status | Verify the MCP's underlying service is running |
| Server logs | Ask server admin to check MCP logs |

**Real example:** URL was `https://kirha-mcp-6022.up.railway.app/mcp\t Open menu` — contained a tab character and extra text. Solution: Use clean URL.

---

## Error 500: Server Error

**Error message:**
```json
{
  "error": {
    "message": "An error occurred while processing your request...",
    "type": "server_error",
    "param": null,
    "code": "server_error"
  }
}
```

**Causes:**
1. MCP server crashed
2. Out of credits (Kirha)
3. Backend service unavailable
4. Request timeout

**Solutions:**

| Check | Action |
|-------|--------|
| Credits | For Kirha: verify you have remaining credits |
| Server status | Check if MCP server is running |
| Rate limits | Wait and retry if you hit rate limits |
| Request size | Reduce context if sending too much data |

**Real example:** Kirha returning 500 because account had no credits. Solution: Add credits to Kirha account.

---

## "MCPConnection URL already exists for this agent"

**Cause:** You're trying to add an MCP that's already configured for this agent.

**Solutions:**
1. Check **Manage MCPs** — the connection might already exist
2. If you need to update it, delete the existing connection first
3. Use a different name if you need multiple connections to the same URL

---

## "Invalid tool schema"

**Error message:**
```json
{
  "message": "Invalid tool schema for: lpltools.hubspot-search-objects",
  "type": "invalid_request_error",
  "param": "tools",
  "code": null
}
```

**Cause:** The MCP server returns a tool definition that the LLM doesn't understand.

**Solutions:**
1. Contact MCP server administrator to fix the schema
2. Try a different LLM (some are more tolerant of schema variations)
3. If it's your MCP, ensure inputSchema follows JSON Schema spec

---

## Agent Not Using MCP Tools

**Symptoms:**
- MCP is configured
- No errors in logs
- Agent just doesn't call the tools

**Causes:**
1. LLM doesn't support MCP (Mistral, Deepseek, Grok)
2. Agent's LLM configuration doesn't include MCP
3. Prompt doesn't encourage tool usage
4. Agent has tool information in memory that it was unavailable

**Solutions:**

| Check | Action |
|-------|--------|
| LLM compatibility | Ensure using OpenAI, Anthropic, or Bedrock |
| Fresh context | Try a new conversation (agent may have cached "unavailable" status) |
| Prompt tuning | Add instructions about available tools in system prompt |

**Real example:** Agent kept saying "mcplpl is indisponible" because it remembered a previous error. Solution: Start new conversation.

---

## URL Issues

### Missing HTTPS

**Error:** Configuration rejected

**Solution:** MCP URLs must use HTTPS. Change `http://` to `https://`

### Missing/Extra Trailing Slash

**Symptoms:** 404 or connection errors

| MCP | Correct URL | Wrong URL |
|-----|-------------|-----------|
| Markitdown | `https://...railway.app/mcp/` | `https://...railway.app/mcp` |
| Kirha | `https://...railway.app/mcp` | `https://...railway.app/mcp/` |

**Solution:** Try adding or removing the trailing slash.

### URL with Hidden Characters

**Symptoms:** 424 or connection errors

**Check for:**
- Tab characters (`\t`)
- "Open menu" or other UI text accidentally copied
- Leading/trailing spaces

**Solution:** Manually retype the URL or use a text editor to inspect for hidden characters.

---

## Anthropic-Specific Issues

### Context Too Long

**Error:** Request fails when using Anthropic with MCP

**Cause:** Anthropic has stricter token limits than OpenAI.

**Solution:** 
1. Reduce `MAX_TOKENS_PER_REQUEST` environment variable to 50,000 (from default 100,000)
2. Or switch to a model with larger context window

---

## Testing MCP Connections

### Using cURL

```bash
# Test tools/list
curl -X POST 'https://your-mcp-server.com/mcp' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer YOUR_API_KEY' \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "params": {},
    "id": 1
  }'
```

### Using Postman

1. Create a POST request
2. URL: Your MCP endpoint
3. Headers:
   - `Content-Type: application/json`
   - `Accept: application/json`
   - `Authorization: Bearer YOUR_KEY` (or your auth header)
4. Body (raw JSON):
```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "params": {},
  "id": 1
}
```

### Using MCP Inspector

```bash
npx @modelcontextprotocol/inspector https://your-mcp-server.com/mcp
```

This opens a visual interface to explore available tools.

---

## Checking Logs

If you have access to orchestrator logs, look for:

```
# Successful MCP call
mcp_list_tools → response with tools array

# Failed MCP call
external_connector_error → check the error message
```

Ask your 6022 administrator if you need help accessing logs.

---

## Getting Help

If you've tried the above solutions and still have issues:

1. **Document the error** — Full error message, JSON if available
2. **Note the MCP** — Which MCP server, what URL
3. **Note the LLM** — Which provider and model
4. **Test independently** — Does the MCP work via cURL/Postman?
5. **Contact support** — Provide all the above information