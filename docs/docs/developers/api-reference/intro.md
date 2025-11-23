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
  --api-base-url https://orchestrator-api-polygon-orchestrator-myproject.up.railway.app/
```

| Network         | Chain ID |
| --------------- | -------- |
| Polygon Mainnet | `137`    |
| Amoy Testnet    | `80002`  |

### JWT Properties

| Property | Value                            |
| -------- | -------------------------------- |
| Validity | 1 month                          |
| Scope    | Per owner wallet (not per agent) |
| Signing  | Owner's private key              |

Include the token in the `Authorization` header (`Authorization: Bearer`).

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

| Status | Meaning                                        |
| ------ | ---------------------------------------------- |
| `200`  | Success (message appended or resource created) |
| `400`  | Bad request (validation failure, duplicates)   |
| `401`  | Unauthorized (JWT missing or invalid)          |
| `404`  | Not found (context key missing)                |
| `500`  | Internal server error                          |

## Python Client

```python
import requests

class AgentMemoryAPI:
    def __init__(self, base_url: str, agent_address: str, agent_token_id: str, jwt: str):
        self.base_url = base_url
        self.agent_address = agent_address
        self.agent_token_id = agent_token_id
        self.headers = {
            "Authorization": f"Bearer {jwt}",
            "Content-Type": "application/json"
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
            return False  # Key already present
        response.raise_for_status()

    def add_memory(self, key: str, message: str, role: str = "agent") -> None:
        url = f"{self.base_url}/agent/{self.agent_address}/{self.agent_token_id}/context/{key}/message"
        response = requests.post(url, headers=self.headers, json={"message": message, "role": role})
        response.raise_for_status()

    def hydrate(self, key: str, conversation: list) -> None:
        self.create_key(key)
        for msg in conversation:
            self.add_memory(key, msg["message"], msg["role"])

# Usage
api = AgentMemoryAPI(
    base_url="https://orchestrator-api-polygon-orchestrator-myproject.up.railway.app",
    agent_address="0x0bd248B276b7Ba8E9E4A94df8dd73CFe9dD2E18a",
    agent_token_id="2",
    jwt="your_jwt_token"
)

keys = api.list_keys()
print("Existing keys:", keys)

api.create_key("C99999")
api.hydrate("C99999", [
    {"role": "user", "message": "Remember: Customer C99999 is a VIP"},
    {"role": "agent", "message": "Noted. C99999 is flagged as VIP."}
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

  async addMemory(key, message, role = "agent") {
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
