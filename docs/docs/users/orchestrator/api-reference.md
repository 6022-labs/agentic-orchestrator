---
sidebar_position: 6
---

# API Reference

The orchestrator provides a REST API for programmatic access to agent management. This is useful for:

- Automated memory updates
- Batch context creation
- Integration with external systems
- CI/CD pipelines

## Base URLs

| Environment | API Base URL |
|-------------|--------------|
| **Production (Polygon)** | `https://orchestrator-api-polygon.up.railway.app` |
| **Testnet (Amoy)** | `https://orchestrator-api-amoy.up.railway.app` |

## Authentication

All API requests require a JWT token signed by the agent owner's wallet.

### Generating a JWT Token

Use the 6022 JWT generator script:

```bash
git clone https://github.com/6022protocol/py-generate-jwt
cd py-generate-jwt

python generate_jwt.py \
  --private-key YOUR_PRIVATE_KEY \
  --chain-id 137 \
  --api-base-url https://orchestrator-api-polygon.up.railway.app/
```

| Parameter | Description |
|-----------|-------------|
| `--private-key` | Wallet private key (owner of the agent) |
| `--chain-id` | `137` for Polygon, `80002` for Amoy |
| `--api-base-url` | Orchestrator API URL |

### Using the Token

Include the JWT in the `Authorization` header:

```bash
curl -X GET \
  https://orchestrator-api-polygon.up.railway.app/agent/1/contexts/keys \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

:::warning Security
Use a dedicated wallet for API access. Never expose your main wallet's private key in scripts.
:::

---

## Context Management Endpoints

### List Context Keys

Get all context keys for an agent.

```http
GET /agent/:tokenId/contexts/keys
```

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `tokenId` | path | Agent's token ID in the collection |

**Response:**
```json
{
  "keys": [
    "POL-2024-123456",
    "POL-2024-789012",
    "INV-2024-001"
  ]
}
```

---

### Create Context

Create a new context key for an agent.

```http
POST /agent/:tokenId/context
```

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `tokenId` | path | Agent's token ID |

**Body:**
```json
{
  "key": "POL-2024-123456"
}
```

**Response:**
```json
{
  "success": true,
  "key": "POL-2024-123456"
}
```

**Errors:**
| Code | Description |
|------|-------------|
| 400 | Context key already exists |
| 401 | Unauthorized (invalid JWT) |
| 403 | Not owner of agent |

---

### Add Memory to Context

Append a message to an agent's context memory.

```http
POST /agent/:tokenId/context/:key/memory
```

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `tokenId` | path | Agent's token ID |
| `key` | path | Context key |

**Body:**
```json
{
  "role": "user",
  "content": "Customer called about policy renewal"
}
```

**Roles:**
| Role | Description |
|------|-------------|
| `user` | External input (customer, user) |
| `assistant` | Agent's response |
| `system` | System message |

**Response:**
```json
{
  "success": true,
  "messageId": "msg_abc123"
}
```

---

### Get Context History

Retrieve conversation history for a context.

```http
GET /agent/:tokenId/context/:key/history
```

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `tokenId` | path | Agent's token ID |
| `key` | path | Context key |

**Response:**
```json
{
  "key": "POL-2024-123456",
  "messages": [
    {
      "role": "user",
      "content": "What is the status of my claim?",
      "timestamp": "2024-01-15T10:30:00Z"
    },
    {
      "role": "assistant",
      "content": "Your claim is currently under review...",
      "timestamp": "2024-01-15T10:30:05Z"
    }
  ]
}
```

---

## Memory Hydration

For bulk memory updates, use the hydration endpoint.

### Hydrate Agent Memory

```http
POST /agent/:tokenId/hydrate
```

**Body:**
```json
{
  "contexts": [
    {
      "key": "POL-2024-001",
      "messages": [
        {"role": "system", "content": "Policy created on 2024-01-01"},
        {"role": "user", "content": "Initial claim filed"}
      ]
    },
    {
      "key": "POL-2024-002",
      "messages": [
        {"role": "system", "content": "Policy created on 2024-02-01"}
      ]
    }
  ]
}
```

**Use cases:**
- Initial data migration
- Restoring from backup
- Syncing from external CRM

---

## Code Examples

### Python

```python
import requests

BASE_URL = "https://orchestrator-api-polygon.up.railway.app"
TOKEN = "your_jwt_token"
AGENT_ID = 1

headers = {
    "Authorization": f"Bearer {TOKEN}",
    "Content-Type": "application/json"
}

# Create a context
response = requests.post(
    f"{BASE_URL}/agent/{AGENT_ID}/context",
    headers=headers,
    json={"key": "POL-2024-123456"}
)
print(response.json())

# Add memory
response = requests.post(
    f"{BASE_URL}/agent/{AGENT_ID}/context/POL-2024-123456/memory",
    headers=headers,
    json={
        "role": "user",
        "content": "Customer inquiry about coverage"
    }
)
print(response.json())
```

### JavaScript / Node.js

```javascript
const axios = require('axios');

const BASE_URL = 'https://orchestrator-api-polygon.up.railway.app';
const TOKEN = 'your_jwt_token';
const AGENT_ID = 1;

const api = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Authorization': `Bearer ${TOKEN}`,
    'Content-Type': 'application/json'
  }
});

// Create a context
async function createContext(key) {
  const response = await api.post(`/agent/${AGENT_ID}/context`, { key });
  return response.data;
}

// Add memory
async function addMemory(key, role, content) {
  const response = await api.post(
    `/agent/${AGENT_ID}/context/${key}/memory`,
    { role, content }
  );
  return response.data;
}

// Usage
createContext('POL-2024-123456')
  .then(() => addMemory('POL-2024-123456', 'user', 'Hello'))
  .then(console.log);
```

### cURL

```bash
# Set variables
TOKEN="your_jwt_token"
BASE="https://orchestrator-api-polygon.up.railway.app"
AGENT=1

# List context keys
curl -X GET "$BASE/agent/$AGENT/contexts/keys" \
  -H "Authorization: Bearer $TOKEN"

# Create context
curl -X POST "$BASE/agent/$AGENT/context" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"key": "POL-2024-123456"}'

# Add memory
curl -X POST "$BASE/agent/$AGENT/context/POL-2024-123456/memory" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "user", "content": "Customer message"}'
```

---

## Error Handling

### Common Error Codes

| Code | Meaning | Solution |
|------|---------|----------|
| 400 | Bad Request | Check request body format |
| 401 | Unauthorized | JWT invalid or expired |
| 403 | Forbidden | Not owner of agent |
| 404 | Not Found | Agent or context doesn't exist |
| 409 | Conflict | Context key already exists |
| 500 | Server Error | Contact support |

### Error Response Format

```json
{
  "error": true,
  "message": "Context key already exists",
  "code": "CONTEXT_EXISTS"
}
```

---

## Rate Limits

| Endpoint | Limit |
|----------|-------|
| Context operations | 100 requests/minute |
| Memory operations | 500 requests/minute |
| Hydration | 10 requests/minute |

Exceeding limits returns `429 Too Many Requests`.