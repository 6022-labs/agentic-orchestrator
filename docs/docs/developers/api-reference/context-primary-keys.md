---
label: Context Primary Keys API
title: Programmatic Context Primary Keys
sidebar_position: 2
---

# Programmatic Context Primary Keys

This guide covers everything you need to provision, query, and manage context primary keys through the Orchestrator API. Use it alongside the [user-facing overview](../../users/agent/agent-knowledge/context-primary-keys.md) when you need to automate onboarding, hydration pipelines, or large-scale migrations.

## Prerequisites

To call these endpoints you need:

- A **[valid JWT](./intro.md#authentication)** signed with the owner wallet of the agent you manage
- The agent's **collection address** and **token ID**
- A chosen **context primary key value** that matches your domain format

Generate a JWT with the official helper project:

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

Tokens are currently valid for one month (planned to decrease to one day). Rotate credentials whenever team members change.

## REST Endpoints

All endpoints live under the orchestrator API base URL:

```
https://orchestrator-api-{network}-orchestrator-{project}.up.railway.app/
```

Swagger documentation: `{base_url}/swagger/index.html`

### List Context Primary Keys

```http
GET /agent/{agentCollectionAddress}/{agentTokenId}/contexts/keys
Authorization: Bearer
```

Path parameters:

| Parameter                | Type   | Description                          |
| ------------------------ | ------ | ------------------------------------ |
| `agentCollectionAddress` | string | Agent NFT collection address (0x...) |
| `agentTokenId`           | string | Agent's token ID within collection   |

Query parameters:

| Parameter         | Type   | Description                                        | Default |
| ----------------- | ------ | -------------------------------------------------- | ------- |
| `page`            | uint   | Page number (1-indexed).                           | `1`     |
| `elementsPerPage` | uint   | Page size, clamped between 1 and 100.              | `10`    |
| `search`          | string | Optional substring filter applied to context keys. | —       |
| `sortBy`          | string | Sort field. Only `key` is currently supported.     | `key`   |
| `sortDirection`   | string | Sort order (`asc` or `desc`).                      | `asc`   |

Successful response (`200 OK`):

```json
{
  "page": 1,
  "elementsPerPage": 10,
  "totalElements": 3,
  "contextKeys": ["C45653", "C78901", "C12345"]
}
```

### Create Context Primary Key

```http
POST /agent/{agentCollectionAddress}/{agentTokenId}/context
Authorization: Bearer
Content-Type: application/json
```

Body:

```json
{
  "key": "C45653"
}
```

- `200 OK`: key provisioned (empty response body)
- `400 Bad Request`: validation failure (bad format, not owner, agent not internal/expert)
- `500 Internal Server Error`: downstream storage error (for example, duplicate keys at persistence time)

### Delete Context Primary Keys

Deletion is intentionally not exposed through the public API. Use database tooling in exceptional cases only (see _Administrative access_ below).

## Bulk Provisioning Example (Python)

```python
import requests

BASE_URL = "https://orchestrator-api-polygon-orchestrator-myproject.up.railway.app"
AGENT_ADDRESS = "0x0bd248B276b7Ba8E9E4A94df8dd73CFe9dD2E18a"
AGENT_TOKEN_ID = "2"
JWT = "your_jwt_token_here"

headers = {
    "Authorization": f"Bearer {JWT}",
    "Content-Type": "application/json"
}

keys_to_create = ["C45653", "C78901", "C12345", "C67890"]

for key in keys_to_create:
    response = requests.post(
        f"{BASE_URL}/agent/{AGENT_ADDRESS}/{AGENT_TOKEN_ID}/context",
        headers=headers,
        json={"key": key}
    )

  if response.status_code == 200:
        print(f"✅ Created: {key}")
  elif response.status_code == 400 and "already exists" in response.text:
    print(f"⏭️  Skipped (exists): {key}")
    else:
        print(f"❌ Error for {key}: {response.text}")
```

## Managing Keys Programmatically

| Action         | Endpoint                                                       | Notes                                             |
| -------------- | -------------------------------------------------------------- | ------------------------------------------------- |
| List keys      | `GET /agent/{address}/{tokenId}/contexts/keys`                 | Returns paginated envelope with `contextKeys`     |
| Create key     | `POST /agent/{address}/{tokenId}/context`                      | Idempotent thanks to duplicate check              |
| Hydrate memory | `POST /agent/{address}/{tokenId}/context/{contextKey}/message` | See [Memory Hydration API](./memory-hydration.md) |

When building automations, always validate the key format locally before calling the API to avoid noisy 400 responses.

## Administrative Access (Direct DB)

Direct database queries should only be used for break-glass situations. Typical queries:

```sql
-- View all contexts for an agent
SELECT * FROM agent_contexts WHERE agent_id = 'agent-uuid';

-- Count context primary keys
SELECT COUNT(*) FROM agent_contexts WHERE agent_id = 'agent-uuid';
```

Coordinate with the platform team before running manual SQL to avoid violating tenant isolation guarantees.

## Troubleshooting

| Status | Error                                     | Cause                             | Fix                                            |
| ------ | ----------------------------------------- | --------------------------------- | ---------------------------------------------- |
| `400`  | `agent context with key X already exists` | Key previously created            | Treat as success, skip duplicates              |
| `400`  | `You are not the owner of this agent`     | JWT signed with wrong wallet      | Regenerate JWT with the current owner          |
| `401`  | `Unauthorized`                            | Missing or invalid JWT            | Refresh token, confirm header formatting       |
| `404`  | `Context not found`                       | Key missing when hydrating memory | Create key before pushing conversations        |
| `500`  | `Internal Server Error`                   | Database or sync issue            | Check orchestrator logs, open ticket if needed |

For UI-driven workflows, see the [user documentation](../../users/agent/agent-knowledge/context-primary-keys.md#creating-context-primary-keys).
