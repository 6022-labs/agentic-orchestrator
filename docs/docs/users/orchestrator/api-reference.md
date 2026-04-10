---
sidebar_position: 6
---

# API Reference

The orchestrator exposes a REST API for programmatic access to agent management. Use it for:

- Memory hydration (injecting messages into an agent's context)
- Batch context creation
- Integration with external systems
- CI/CD pipelines

All endpoints live behind two path parameters: the **agent collection address** (the NFT contract) and the **agent token ID** (the agent's index in that collection). You will find both in the orchestrator dashboard for any agent you own.

## Base URLs

| Environment              | API Base URL                                |
| ------------------------ | ------------------------------------------- |
| **Production (Polygon)** | `https://<your-orchestrator-api-host>`      |
| **Testnet (Amoy)**       | `https://<your-orchestrator-api-host-amoy>` |

Replace the placeholders with the deployment URLs your team operates.

## Authentication

Every request requires a JWT issued by the SIWE flow described in **[Web3 Authentication](./authentication.md)**. The flow is:

1. `POST /auth/jwt/init` → returns a nonce + SIWE message
2. Sign the message locally with the wallet that owns the agent
3. `POST /auth/jwt/finish` → returns a JWT
4. Send `Authorization: Bearer <token>` on every subsequent call

You can run the flow manually (see Authentication for examples) or use the helper script:

```bash
git clone https://github.com/6022protocol/py-generate-jwt
cd py-generate-jwt

python generate_jwt.py \
  --private-key 0x<your private key> \
  --chain-id 137 \
  --api-base-url https://<your-orchestrator-api-host>/
```

| Parameter        | Description                               |
| ---------------- | ----------------------------------------- |
| `--private-key`  | Wallet private key (owner of the agent)   |
| `--chain-id`     | `137` for Polygon, `80002` for Amoy       |
| `--api-base-url` | Orchestrator API URL, with trailing slash |

Tokens are bound to their `chainId` — a Polygon JWT is rejected by an Amoy orchestrator and vice versa.

```bash
curl -X GET \
  "https://<your-orchestrator-api-host>/agent/0xCOLLECTION/1/contexts/keys" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

:::warning Security
Use a dedicated wallet for API access. Never expose your main wallet's private key in scripts, CI variables, or `.env` files you don't control.
:::

---

## Context Management Endpoints

All paths take **two** path parameters: `:agentCollectionAddress` (the contract) and `:agentCollectionTokenId` (the agent's token ID inside that contract).

### List context keys

Get all context keys for an agent.

```http
GET /agent/:agentCollectionAddress/:agentCollectionTokenId/contexts/keys
```

**Path parameters**

| Name                     | Description                               |
| ------------------------ | ----------------------------------------- |
| `agentCollectionAddress` | NFT collection contract address (`0x...`) |
| `agentCollectionTokenId` | Agent's token ID inside the collection    |

**Response**

```json
{
  "keys": ["POL-2024-123456", "POL-2024-789012", "INV-2024-001"]
}
```

Internal context keys (anything under the `conversation-context-discovery/` prefix) are filtered out of the response.

---

### Create a context

Create a new context key for an agent.

```http
POST /agent/:agentCollectionAddress/:agentCollectionTokenId/context
```

**Body**

```json
{
  "key": "POL-2024-123456"
}
```

**Response**

```json
{
  "success": true,
  "key": "POL-2024-123456"
}
```

**Errors**

| Code  | Cause                                                                                 |
| ----- | ------------------------------------------------------------------------------------- |
| `400` | The key is empty or starts with the reserved `conversation-context-discovery/` prefix |
| `400` | A context with the same key already exists for this agent                             |
| `401` | Missing or invalid JWT                                                                |
| `403` | The wallet in your JWT does not own this agent                                        |

---

### Append a message to a context (memory hydration)

Add a message to an agent's context memory. **This is the endpoint to use for hydration** — there is no separate `/hydrate` endpoint. Call it once per message.

```http
POST /agent/:agentCollectionAddress/:agentCollectionTokenId/context/:contextKey/message
```

**Path parameters**

| Name                     | Description                                       |
| ------------------------ | ------------------------------------------------- |
| `agentCollectionAddress` | NFT collection contract address                   |
| `agentCollectionTokenId` | Agent's token ID                                  |
| `contextKey`             | The context key to append to (must already exist) |

**Body**

```json
{
  "message": "Customer called about policy renewal",
  "role": "user"
}
```

`message` and the context key must each be at least 2 characters. The role enum is enforced server-side.

**Roles**

| Role        | Description                         |
| ----------- | ----------------------------------- |
| `user`      | External input (customer, end user) |
| `assistant` | Agent's response                    |
| `system`    | System message                      |

**Response**

```
200 OK
```

If the context key does not exist, you'll get `400` with an error message — create the context first via the endpoint above.

---

## Code Examples

### Python

```python
import requests

BASE_URL = "https://<your-orchestrator-api-host>"
TOKEN = "your_jwt_token"
COLLECTION = "0xCOLLECTION"
TOKEN_ID = 1

headers = {
    "Authorization": f"Bearer {TOKEN}",
    "Content-Type": "application/json",
}

# Create a context
requests.post(
    f"{BASE_URL}/agent/{COLLECTION}/{TOKEN_ID}/context",
    headers=headers,
    json={"key": "POL-2024-123456"},
)

# Hydrate it with a message
requests.post(
    f"{BASE_URL}/agent/{COLLECTION}/{TOKEN_ID}/context/POL-2024-123456/message",
    headers=headers,
    json={
        "message": "Customer inquiry about coverage",
        "role": "user",
    },
)
```

### JavaScript / Node.js

```javascript
const axios = require("axios");

const BASE_URL = "https://<your-orchestrator-api-host>";
const TOKEN = "your_jwt_token";
const COLLECTION = "0xCOLLECTION";
const TOKEN_ID = 1;

const api = axios.create({
  baseURL: BASE_URL,
  headers: {
    Authorization: `Bearer ${TOKEN}`,
    "Content-Type": "application/json",
  },
});

await api.post(`/agent/${COLLECTION}/${TOKEN_ID}/context`, {
  key: "POL-2024-123456",
});

await api.post(
  `/agent/${COLLECTION}/${TOKEN_ID}/context/POL-2024-123456/message`,
  { message: "Hello", role: "user" },
);
```

### cURL

```bash
TOKEN="your_jwt_token"
BASE="https://<your-orchestrator-api-host>"
COLLECTION="0xCOLLECTION"
TOKEN_ID=1

# List context keys
curl -X GET "$BASE/agent/$COLLECTION/$TOKEN_ID/contexts/keys" \
  -H "Authorization: Bearer $TOKEN"

# Create context
curl -X POST "$BASE/agent/$COLLECTION/$TOKEN_ID/context" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"key": "POL-2024-123456"}'

# Append a message (hydration)
curl -X POST "$BASE/agent/$COLLECTION/$TOKEN_ID/context/POL-2024-123456/message" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"message": "Customer message", "role": "user"}'
```

---

## Other endpoints

The orchestrator also exposes administrative endpoints used by the dashboard UI:

| Concern                             | Endpoints                                                                              |
| ----------------------------------- | -------------------------------------------------------------------------------------- |
| LLM settings                        | `GET/POST /agent/:agentCollectionAddress/:agentCollectionTokenId/llm/settings`         |
| Database setup (two-step)           | `POST .../database/settings/first-step`, `POST .../database/settings/second-step`      |
| Prompt templates                    | `GET/POST /agent/:agentCollectionAddress/:agentCollectionTokenId/prompt-template/:key` |
| MCP connections                     | `GET .../mcp/connections`, `POST .../mcp/connection`, `DELETE .../mcp/connection/:id`  |
| Wallet verification (Slack pairing) | `POST /wallet/verify`                                                                  |

These are used by the dashboard and the Slack bridge — most users do not need to call them directly. If you need the full request/response shapes, read the gateway source under `gateway/src/gateway_mvc/` or the developer reference.

---

## Error Handling

### Common error codes

| Code  | Meaning                                                 | Solution                                      |
| ----- | ------------------------------------------------------- | --------------------------------------------- |
| `400` | Bad request — body schema, validation, or duplicate key | Check the response message                    |
| `401` | Unauthorized — JWT invalid, expired, or missing         | Re-issue a JWT (see Authentication)           |
| `403` | Forbidden — your wallet does not own this agent         | Confirm NFT ownership in the dashboard        |
| `404` | Agent or context does not exist                         | Check the collection address / token ID / key |

### Error response format

The orchestrator currently surfaces errors as plain text or short JSON depending on the controller. Future versions may standardise on a single shape — do not rely on a specific JSON structure beyond `Authorization`-related codes.
