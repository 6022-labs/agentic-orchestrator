---
title: API Reference Overview
sidebar_position: 1
---

# API Reference Overview

The 6022 Orchestrator exposes a REST API for programmatic access to agent management. Use it for:

- Memory hydration (injecting messages into an agent's context)
- Batch context creation
- Integration with external systems
- CI/CD pipelines

All resource endpoints live behind two path parameters: the **agent collection address** (the NFT contract) and the **agent token ID** (the agent's index in that collection). You'll find both in the orchestrator dashboard for any agent you own.

## Base URLs

| Environment              | API Base URL                                |
| ------------------------ | ------------------------------------------- |
| **Production (Polygon)** | `https://<your-orchestrator-api-host>`      |
| **Testnet (Amoy)**       | `https://<your-orchestrator-api-host-amoy>` |

Replace the placeholders with the deployment URLs your team operates.

## Authentication

Every request requires a JWT obtained by signing a SIWE message with the wallet that owns the agent. The full two-step flow, JWT claim structure, and helper-script usage is documented on a dedicated page:

→ **[Authentication](./authentication.md)** — the SIWE/JWT flow and how to generate a token

Short version:

1. `POST /auth/jwt/init` → returns a nonce + SIWE message
2. Sign the message locally with the wallet that owns the agent
3. `POST /auth/jwt/finish` → returns a JWT
4. Send `Authorization: Bearer <token>` on every subsequent call

Tokens are bound to their `chainId` — a Polygon JWT is rejected by an Amoy orchestrator and vice versa. Tokens expire after `JWT__LIVING_TIME` seconds (deployment-configured); repeat the SIWE flow to get a new one.

```bash
curl -X GET \
  "https://<your-orchestrator-api-host>/agent/0xCOLLECTION/1/contexts/keys" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

:::warning Security
Use a dedicated wallet for API access. Never expose your main wallet's private key in scripts, CI variables, or `.env` files you don't control.
:::

---

## Dedicated reference pages

Rather than document every endpoint on this page, the detailed reference lives in topic-specific pages:

| Page | What it covers |
|---|---|
| **[Authentication](./authentication.md)** | SIWE nonce flow, JWT claims, `JWT__HASH` / `JWT__LIVING_TIME`, helper script, Python sample |
| **[Context Primary Keys API](./context-primary-keys.md)** | `GET /contexts/keys`, `POST /context`, pagination/search/sort parameters, bulk provisioning, administrative DB access |
| **[Memory Hydration API](./memory-hydration.md)** | `POST /context/:key/message`, role enum (`user` / `system` / `assistant`), hydration patterns, bulk templates |

Start with **Authentication**, then pick whichever resource page matches what you're integrating.

---

## Context Management — quick reference

For convenience, here's the single-page cheat sheet for the two most common endpoint groups. The dedicated reference pages above have the full parameter documentation.

### List context keys

```http
GET /agent/:agentCollectionAddress/:agentCollectionTokenId/contexts/keys
Authorization: Bearer <token>
```

Internal context keys (anything under the `conversation-context-discovery/` prefix) are filtered out of the response. Returns a paginated envelope — see [Context Primary Keys API](./context-primary-keys.md) for query parameters.

### Create a context

```http
POST /agent/:agentCollectionAddress/:agentCollectionTokenId/context
Authorization: Bearer <token>
Content-Type: application/json

{ "key": "POL-2024-123456" }
```

### Append a message to a context (memory hydration)

```http
POST /agent/:agentCollectionAddress/:agentCollectionTokenId/context/:contextKey/message
Authorization: Bearer <token>
Content-Type: application/json

{
  "message": "Customer called about policy renewal",
  "role": "user"
}
```

`role` must be one of `user`, `system`, or `assistant`. `message` and `contextKey` must each be at least 2 characters. The context key must already exist; create it first via the endpoint above.

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

### Python client class

For longer-running integrations, a small wrapper class keeps things tidy:

```python
import requests


class AgentMemoryAPI:
    """Minimal client for the orchestrator's context/memory endpoints."""

    def __init__(self, base_url: str, agent_address: str, agent_token_id: str, jwt: str):
        self.base_url = base_url
        self.agent_address = agent_address
        self.agent_token_id = agent_token_id
        self.headers = {
            "Authorization": f"Bearer {jwt}",
            "Content-Type": "application/json",
        }

    def list_keys(self) -> list:
        url = f"{self.base_url}/agent/{self.agent_address}/{self.agent_token_id}/contexts/keys"
        response = requests.get(url, headers=self.headers)
        response.raise_for_status()
        payload = response.json()
        return payload.get("contextKeys", [])

    def create_key(self, key: str) -> bool:
        url = f"{self.base_url}/agent/{self.agent_address}/{self.agent_token_id}/context"
        response = requests.post(url, headers=self.headers, json={"key": key})

        if response.status_code == 200:
            return True
        if response.status_code == 400 and "already exists" in response.text:
            return False  # Treat duplicates as idempotent success
        response.raise_for_status()
        return False

    def add_memory(self, key: str, message: str, role: str = "assistant") -> None:
        """role must be one of: user, system, assistant."""
        url = f"{self.base_url}/agent/{self.agent_address}/{self.agent_token_id}/context/{key}/message"
        response = requests.post(url, headers=self.headers, json={"message": message, "role": role})
        response.raise_for_status()

    def hydrate(self, key: str, conversation: list) -> None:
        self.create_key(key)
        for msg in conversation:
            self.add_memory(key, msg["message"], msg["role"])


# Usage
api = AgentMemoryAPI(
    base_url="https://<your-orchestrator-api-host>",
    agent_address="0xCOLLECTION",
    agent_token_id="2",
    jwt="your_jwt_token",
)

api.create_key("C99999")
api.hydrate("C99999", [
    {"role": "user", "message": "Remember: Customer C99999 is a VIP"},
    {"role": "assistant", "message": "Noted. C99999 is flagged as VIP."},
])
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

// role must be "user", "system", or "assistant"
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

The orchestrator also exposes administrative endpoints used by the dashboard UI and the Slack bridge:

| Concern                             | Endpoints                                                                              |
| ----------------------------------- | -------------------------------------------------------------------------------------- |
| LLM settings                        | `GET/POST /agent/:agentCollectionAddress/:agentCollectionTokenId/llm/settings`         |
| Database setup (two-step)           | `POST .../database/settings/first-step`, `POST .../database/settings/second-step`      |
| Prompt templates                    | `GET/POST /agent/:agentCollectionAddress/:agentCollectionTokenId/prompt-template/:key` |
| MCP connections                     | `GET .../mcp/connections`, `POST .../mcp/connection`, `DELETE .../mcp/connection/:id`  |
| Wallet verification (Slack pairing) | `POST /wallet/verify`                                                                  |

These are consumed internally today — most integrators do not need to call them directly. If you need the full request/response shapes, read the gateway source under `gateway/src/gateway_mvc/`.

---

## Error Handling

### Common response codes

| Status | Meaning                                                                                                                                         |
| ------ | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| `200`  | Success. The orchestrator's gateway controllers always return `200 OK` for successful POSTs — there is **no** `201 Created`, regardless of what older references may say. |
| `400`  | Bad request (validation failure, duplicate key, schema error)                                                                                   |
| `401`  | Unauthorized — JWT missing, invalid, expired, or chain-id mismatch                                                                              |
| `403`  | Forbidden — your wallet does not own this agent                                                                                                 |
| `404`  | Agent or context not found                                                                                                                      |
| `500`  | Internal server error                                                                                                                           |

### Error response format

The orchestrator currently surfaces errors as plain text or short JSON depending on the controller. Future versions may standardise on a single shape — don't rely on a specific JSON structure beyond `Authorization`-related codes.

### Troubleshooting table

| Error                                     | Cause                             | Fix                                            |
| ----------------------------------------- | --------------------------------- | ---------------------------------------------- |
| `agent context with key X already exists` | Key previously created            | Treat as success, skip duplicates              |
| `You are not the owner of this agent`     | JWT signed with wrong wallet      | Regenerate JWT with the current owner          |
| `Unauthorized`                            | Missing or invalid JWT            | Refresh token, confirm `Authorization` header  |
| `Context not found`                       | Key missing when hydrating memory | Create the key first via `POST /context`       |
| `Internal Server Error`                   | Database or downstream sync issue | Check orchestrator logs, open a ticket if needed |
