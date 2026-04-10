---
title: Introduction
sidebar_position: 1
---

# Introduction

The 6022 Orchestrator exposes REST endpoints for managing agents programmatically. This document mirrors the original reference from the user guide and is the canonical location for developers.

## Authentication

All API requests require JWT (JSON Web Token) authentication.

### Generating a JWT

```bash
# Clone the repository
git clone https://github.com/6022protocol/py-generate-jwt

# Install dependencies
cd py-generate-jwt
pip install -r requirements.txt

# Generate token
python generate_jwt.py \
  --chain-id 137 \
  --api-base-url https://<your-orchestrator-api-host>/
```

| Network         | Chain ID |
| --------------- | -------- |
| Polygon Mainnet | `137`    |
| Amoy Testnet    | `80002`  |

### JWT Properties

| Property | Value                                                                        |
| -------- | ---------------------------------------------------------------------------- |
| Validity | Controlled by `JWT__LIVING_TIME` (seconds) on the orchestrator deployment    |
| Scope    | Per owner wallet (not per agent); chain-bound (`cId` claim)                  |
| Signing  | Orchestrator HMAC secret `JWT__HASH`, signed after SIWE message verification |

Include the token in the `Authorization` header: `Authorization: Bearer <token>`.

A token issued for one chain (e.g. Polygon mainnet `137`) is **not** accepted by an orchestrator instance running against a different chain (e.g. Amoy `80002`) — the `cId` claim is verified on every request.

## Endpoints

### List Context Primary Keys

```http
GET /agent/{agentCollectionAddress}/{agentTokenId}/contexts/keys
Authorization: Bearer
```

### Create Context Primary Key

```http
POST /agent/{agentCollectionAddress}/{agentTokenId}/context
Authorization: Bearer
Content-Type: application/json
```

Request body:

```json
{ "key": "C45653" }
```

### Add Message to Agent Memory

```http
POST /agent/{agentCollectionAddress}/{agentTokenId}/context/{contextKey}/message
Authorization: Bearer
Content-Type: application/json
```

Request body:

```json
{
  "message": "Content of the message to add",
  "role": "user"
}
```

### Response Codes

| Status | Meaning                                                            |
| ------ | ------------------------------------------------------------------ |
| `200`  | Success (resource created or message appended). The orchestrator's gateway controllers always return `200 OK` for successful POSTs — there is no `201 Created` |
| `400`  | Bad request (validation failure, duplicate key, etc.)              |
| `401`  | Unauthorized (JWT missing, invalid, or chain-id mismatch)          |
| `403`  | Forbidden (your wallet does not own this agent)                    |
| `404`  | Not found (agent or context key missing)                           |
| `500`  | Internal server error                                              |

## Python Client

`role` accepts `user`, `system`, or `assistant` — the orchestrator validates against the `AgentContextMessageRole` enum and rejects any other value.

```python
import requests


class AgentMemoryAPI:
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
            return False  # Key already present, treat as idempotent success
        response.raise_for_status()
        return False

    def add_memory(self, key: str, message: str, role: str = "assistant") -> None:
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

keys = api.list_keys()
print("Existing keys:", keys)

api.create_key("C99999")
api.hydrate("C99999", [
    {"role": "user", "message": "Remember: Customer C99999 is a VIP"},
    {"role": "assistant", "message": "Noted. C99999 is flagged as VIP."},
])
```

## JavaScript Client

```javascript
const axios = require("axios");

class AgentMemoryAPI {
  constructor(baseUrl, agentAddress, agentTokenId, jwt) {
    this.baseUrl = baseUrl;
    this.agentAddress = agentAddress;
    this.agentTokenId = agentTokenId;
    this.headers = {
      Authorization: `Bearer ${jwt}`,
      "Content-Type": "application/json",
    };
  }

  async listKeys() {
    const url = `${this.baseUrl}/agent/${this.agentAddress}/${this.agentTokenId}/contexts/keys`;
    const response = await axios.get(url, { headers: this.headers });
    return response.data.contextKeys || [];
  }

  async createKey(key) {
    const url = `${this.baseUrl}/agent/${this.agentAddress}/${this.agentTokenId}/context`;
    try {
      await axios.post(url, { key }, { headers: this.headers });
      return true;
    } catch (error) {
      if (error.response?.data?.message?.includes("already exists")) {
        return false;
      }
      throw error;
    }
  }

  // role must be "user", "system", or "assistant"
  async addMemory(key, message, role = "assistant") {
    const url = `${this.baseUrl}/agent/${this.agentAddress}/${this.agentTokenId}/context/${key}/message`;
    await axios.post(url, { message, role }, { headers: this.headers });
  }
}

module.exports = AgentMemoryAPI;
```

## Troubleshooting

| Error                                     | Cause                        | Solution                                     |
| ----------------------------------------- | ---------------------------- | -------------------------------------------- |
| `agent context with key X already exists` | Key previously created       | Skip or reuse the existing key               |
| `You are not the owner of this agent`     | JWT signed with wrong wallet | Regenerate JWT with owner's private key      |
| `Unauthorized`                            | Missing or invalid JWT       | Check the `Authorization` header             |
| `Context not found`                       | Context key missing          | Create the key first                         |
| `Internal Server Error`                   | Database or sync issue       | Inspect orchestrator logs or contact support |
