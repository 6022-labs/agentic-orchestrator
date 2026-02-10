---
sidebar_position: 5
---

# Building Custom MCPs

When no existing MCP fits your needs, you can build your own. This guide covers the basics of creating an MCP server compatible with 6022 Agentic Teams.

## When to Build Custom

Build a custom MCP when you need to:

- Connect to proprietary internal systems
- Expose company-specific knowledge bases
- Wrap APIs that don't have MCP support
- Combine multiple data sources into one interface

## MCP Architecture

A custom MCP server must implement the JSON-RPC 2.0 protocol with two key methods:

```
┌─────────────────────────────────────────────────────┐
│                   MCP SERVER                        │
├─────────────────────────────────────────────────────┤
│                                                     │
│   POST /mcp                                         │
│                                                     │
│   ┌─────────────┐      ┌─────────────┐             │
│   │ tools/list  │      │ tools/call  │             │
│   │             │      │             │             │
│   │ Returns     │      │ Executes    │             │
│   │ available   │      │ a specific  │             │
│   │ tools       │      │ tool        │             │
│   └─────────────┘      └─────────────┘             │
│                                                     │
└─────────────────────────────────────────────────────┘
```

## Minimal Implementation

### 1. Handle `tools/list`

This method tells the LLM what tools are available:

```python
def handle_tools_list():
    return {
        "tools": [
            {
                "name": "search_documents",
                "description": "Search internal documents by keyword",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "query": {
                            "type": "string",
                            "description": "Search query"
                        },
                        "limit": {
                            "type": "integer",
                            "description": "Max results (default 10)"
                        }
                    },
                    "required": ["query"]
                }
            },
            {
                "name": "get_document",
                "description": "Retrieve a specific document by ID",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "document_id": {
                            "type": "string",
                            "description": "Document identifier"
                        }
                    },
                    "required": ["document_id"]
                }
            }
        ]
    }
```

### 2. Handle `tools/call`

This method executes a tool when the LLM requests it:

```python
def handle_tools_call(tool_name, arguments):
    if tool_name == "search_documents":
        results = search_your_database(
            query=arguments.get("query"),
            limit=arguments.get("limit", 10)
        )
        return {"content": [{"type": "text", "text": json.dumps(results)}]}
    
    elif tool_name == "get_document":
        doc = fetch_document(arguments.get("document_id"))
        return {"content": [{"type": "text", "text": doc}]}
    
    else:
        raise ValueError(f"Unknown tool: {tool_name}")
```

## Full Python Example (Flask)

```python
from flask import Flask, request, jsonify
import json

app = Flask(__name__)

# Your tool implementations
def search_policies(query):
    # Replace with actual search logic
    return [
        {"id": "POL-001", "title": "Remote Work Policy", "snippet": "..."},
        {"id": "POL-002", "title": "Travel Policy", "snippet": "..."},
    ]

def get_policy(policy_id):
    # Replace with actual fetch logic
    policies = {
        "POL-001": {"title": "Remote Work Policy", "content": "Full policy text..."},
        "POL-002": {"title": "Travel Policy", "content": "Full policy text..."},
    }
    return policies.get(policy_id, {"error": "Not found"})

# MCP endpoint
@app.route('/mcp', methods=['POST'])
def mcp_handler():
    data = request.json
    method = data.get('method')
    params = data.get('params', {})
    request_id = data.get('id')
    
    try:
        if method == 'tools/list':
            result = {
                "tools": [
                    {
                        "name": "search_policies",
                        "description": "Search company policies by keyword",
                        "inputSchema": {
                            "type": "object",
                            "properties": {
                                "query": {"type": "string", "description": "Search terms"}
                            },
                            "required": ["query"]
                        }
                    },
                    {
                        "name": "get_policy",
                        "description": "Get full text of a specific policy",
                        "inputSchema": {
                            "type": "object",
                            "properties": {
                                "policy_id": {"type": "string", "description": "Policy ID"}
                            },
                            "required": ["policy_id"]
                        }
                    }
                ]
            }
        
        elif method == 'tools/call':
            tool_name = params.get('name')
            arguments = params.get('arguments', {})
            
            if tool_name == 'search_policies':
                data = search_policies(arguments.get('query'))
                result = {"content": [{"type": "text", "text": json.dumps(data)}]}
            
            elif tool_name == 'get_policy':
                data = get_policy(arguments.get('policy_id'))
                result = {"content": [{"type": "text", "text": json.dumps(data)}]}
            
            else:
                raise ValueError(f"Unknown tool: {tool_name}")
        
        else:
            raise ValueError(f"Unknown method: {method}")
        
        return jsonify({
            "jsonrpc": "2.0",
            "result": result,
            "id": request_id
        })
    
    except Exception as e:
        return jsonify({
            "jsonrpc": "2.0",
            "error": {"code": -32000, "message": str(e)},
            "id": request_id
        }), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=3000)
```

## Adding Authentication

Protect your MCP with authentication:

```python
from functools import wraps

def require_api_key(f):
    @wraps(f)
    def decorated(*args, **kwargs):
        api_key = request.headers.get('x-api-key')
        if api_key != 'your-secret-key':
            return jsonify({"error": "Unauthorized"}), 401
        return f(*args, **kwargs)
    return decorated

@app.route('/mcp', methods=['POST'])
@require_api_key
def mcp_handler():
    # ... handler code
```

Then configure in orchestrator:
- Authentication: Custom Header
- Header Key: `x-api-key`
- Header Value: `your-secret-key`

## Deployment Options

### Railway (Recommended)

1. Push your code to GitHub
2. Connect repo to Railway
3. Deploy automatically
4. Use the Railway URL: `https://your-app.up.railway.app/mcp`

### Docker

```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
CMD ["python", "app.py"]
```

### Fly.io

```bash
flyctl launch
flyctl deploy
```

## Tool Design Best Practices

### Good Tool Descriptions

The LLM uses descriptions to decide when to call tools:

```python
# ❌ Bad - vague
{
    "name": "search",
    "description": "Search stuff"
}

# ✅ Good - specific and actionable
{
    "name": "search_customer_tickets",
    "description": "Search customer support tickets by customer email, ticket ID, or keywords. Returns ticket ID, status, subject, and creation date."
}
```

### Clear Input Schemas

Define what parameters each tool accepts:

```python
{
    "name": "get_invoice",
    "description": "Retrieve invoice details by invoice number",
    "inputSchema": {
        "type": "object",
        "properties": {
            "invoice_number": {
                "type": "string",
                "description": "Invoice number (format: INV-YYYY-NNNN)",
                "pattern": "^INV-\\d{4}-\\d{4}$"
            },
            "include_line_items": {
                "type": "boolean",
                "description": "Whether to include line item details",
                "default": True
            }
        },
        "required": ["invoice_number"]
    }
}
```

### Return Structured Data

Return data that the LLM can easily parse:

```python
# ✅ Good - structured
return {
    "content": [{
        "type": "text",
        "text": json.dumps({
            "customer": {"name": "John", "email": "john@example.com"},
            "tickets": [{"id": "T-001", "status": "open"}]
        })
    }]
}
```

## Testing Your MCP

### Local Testing

```bash
# Start your server
python app.py

# Test tools/list
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":1}'

# Test tools/call
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_policies","arguments":{"query":"remote work"}},"id":2}'
```

### MCP Inspector

```bash
npx @modelcontextprotocol/inspector http://localhost:3000/mcp
```

## Resources

- [MCP Specification](https://modelcontextprotocol.io/specification/2025-06-18/server/tools)
- [MCP Python SDK](https://github.com/modelcontextprotocol/python-sdk)
- [MCP TypeScript SDK](https://github.com/modelcontextprotocol/typescript-sdk)
- [Example MCP Servers](https://github.com/modelcontextprotocol/servers)