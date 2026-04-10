---
sidebar_position: 4
---

# MCP Troubleshooting

Common MCP issues and their fixes. The error JSON shapes shown below come from the upstream OpenAI Responses API wrapping the failure — the orchestrator itself returns plain Go error strings for the same conditions.

## 401 Unauthorized

The MCP server rejected your credentials.

- **Bad or expired key** — generate a new one from the provider's dashboard.
- **Wrong auth method** — some servers use `Authorization: Bearer`, others use a custom header like `x-api-key`. Check the [Available MCPs](./available-mcps.md) table for the right method per server.
- **Stray whitespace** — double-check the header value for leading/trailing spaces.
- **Notion 401** — the Notion MCP requires OAuth, which the orchestrator doesn't currently support. Self-host instead.

## 424 Failed Dependency

The MCP server couldn't reach its own backend. Common causes:

- **URL contains hidden characters** — tabs, UI text like `Open menu` that was accidentally copied, leading/trailing spaces. Retype the URL manually.
- **MCP backend is down** — ask the MCP server admin to check logs.

## 500 Server Error

The MCP server itself crashed or ran out of resources. For Kirha specifically this usually means **out of credits** — check your Kirha account balance. Otherwise, verify the server is running and retry after a short wait.

## "MCPConnection URL already exists for this agent"

You already added this MCP to the agent. Go to **⋮ → Manage MCPs**, edit the existing entry, or delete it and re-add.

## "Invalid tool schema"

The MCP server returned a tool definition that doesn't match the JSON Schema spec. Options: contact the server admin to fix the schema, try a more tolerant LLM provider, or if it's your MCP, validate the `inputSchema` against the [JSON Schema spec](https://json-schema.org/).

## Agent configured with MCP but never uses it

- **LLM doesn't support function calling** — check the [LLM compatibility matrix](./intro.md#llm-compatibility). If you're on Mistral/DeepSeek/xAI, the model has to support function calling for the loop to produce tool calls.
- **Cached "unavailable" in memory** — the agent may have concluded earlier in the conversation that the MCP was down and is now repeating itself. Start a new conversation.
- **Prompt doesn't mention the MCP** — add a section to the System prompt telling the agent what tools it has and when to use them. See [Best Practices → Prompting for MCP usage](./best-practice.md#prompting-for-mcp-usage).

## URL issues

| Problem | Fix |
|---|---|
| `http://` URL | The orchestrator rejects non-HTTPS URLs. Use `https://`. |
| Wrong trailing slash | Markitdown requires `/mcp/` (trailing slash). Kirha uses `/mcp` (no trailing slash). Try both. |
| Hidden characters | Retype the URL manually in a plain-text editor to strip invisible whitespace and pasted UI text. |

## Anthropic-specific: context too long

Anthropic counts tokens more aggressively than OpenAI for the same logical content. Lower `TOKEN__MAX_TOKENS_PER_REQUEST` on the orchestrator to e.g. `50000` (default is `100000`). The `TOKEN__` prefix is required — the bare form is silently ignored. See [LLM Troubleshooting](../llm/troubleshooting.md) for the canonical explanation.

## Testing an MCP connection from outside the orchestrator

Before debugging inside the dashboard, verify the MCP server itself is reachable.

### cURL

```bash
curl -X POST 'https://your-mcp-server.com/mcp' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_API_KEY' \
  -d '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":1}'
```

A working server responds with a `result.tools` array. Anything else — 401, 404, 500, connection refused — is the problem you need to fix before touching the orchestrator.

### MCP Inspector (visual)

```bash
npx @modelcontextprotocol/inspector https://your-mcp-server.com/mcp
```

Opens a web UI that lets you browse the server's tools, call them with custom arguments, and see raw responses. Useful when the server says it has tools but the orchestrator isn't picking them up.

## If nothing works

Collect these before asking for help:

1. The full error message (including any JSON body)
2. Which MCP server and URL
3. Which LLM provider and model on the agent
4. Whether the MCP works via cURL / MCP Inspector (i.e. is it an MCP problem or an orchestrator problem?)
