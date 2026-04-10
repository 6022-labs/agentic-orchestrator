---
title: Memory Hydration API
sidebar_position: 3
---

# Memory Hydration API

Use this guide when you need to preload agent memory with synthetic conversations. It complements the conceptual overview in the [user documentation](../../users/memory/memory-hydration.md).

## Prerequisites

- Context primary keys must already exist for the agent you want to hydrate
- A [valid JWT](./intro.md#authentication) signed by the agent owner
- The agent collection address and token ID

## Endpoint

```http
POST /agent/{agentCollectionAddress}/{agentTokenId}/context/{contextKey}/message
Authorization: Bearer
Content-Type: application/json
```

| Parameter                | Description                          |
| ------------------------ | ------------------------------------ |
| `agentCollectionAddress` | Agent NFT collection address         |
| `agentTokenId`           | Agent token ID within the collection |
| `contextKey`             | Target context primary key           |

Body schema:

```json
{
  "message": "Content of the memory to add",
  "role": "user"
}
```

`role` accepts the values defined by the `AgentContextMessageRole` enum: `user` (messages injected as if a human said them), `assistant` (messages representing the agent's response or stored notes), or `system` (system-level instructions). The orchestrator rejects any other value with a 400.

Successful calls return `200 OK`. A missing context key yields `404 Not Found`.

## Hydration Patterns

Reuse these snippets to keep your hydrations consistent across multiple keys.

### Simple Note

```json
{
  "message": "Customer John Doe, contract C45653, VIP status, prefers email contact.",
  "role": "user"
}
```

```json
{
  "message": "Noted: Customer John Doe (C45653) is a VIP customer with email preference.",
  "role": "assistant"
}
```

### Knowledge Transfer

```json
{
  "message": "Note for later the whole content of GDPR chapter 'individual privacy 3.1': ...",
  "role": "user"
}
```

```json
{
  "message": "Yes, duly noted. I have stored the content of GDPR chapter 3.1 regarding individual privacy rights.",
  "role": "assistant"
}
```

### Multi-Turn Conversation

```json
{ "message": "What are the steps for processing a refund request?", "role": "user" }
{ "message": "The refund process has 4 steps: 1) Verify purchase within 30 days, 2) Check item condition, 3) Process in CRM, 4) Confirm with customer.", "role": "assistant" }
{ "message": "What if the item is damaged?", "role": "user" }
{ "message": "For damaged items, skip step 2 and escalate to supervisor before processing. Document damage with photos.", "role": "assistant" }
```

### Persona Establishment

```json
{ "message": "Remember that you should always be empathetic with customers who have complaints.", "role": "user" }
{ "message": "Understood. I will approach all complaints with empathy and understanding, acknowledging the customer's frustration before offering solutions.", "role": "assistant" }
```

## Python Helper

```python
import requests

BASE_URL = "https://<your-orchestrator-api-host>"
AGENT_ADDRESS = "0xCOLLECTION"
AGENT_TOKEN_ID = "2"
CONTEXT_KEY = "Article-15"
JWT = "your_jwt_token_here"

headers = {
    "Authorization": f"Bearer {JWT}",
    "Content-Type": "application/json",
}


def add_memory(message: str, role: str) -> None:
    """role must be one of: user, system, assistant."""
    url = f"{BASE_URL}/agent/{AGENT_ADDRESS}/{AGENT_TOKEN_ID}/context/{CONTEXT_KEY}/message"
    response = requests.post(url, headers=headers, json={"message": message, "role": role})

    if response.status_code == 200:
        print(f"✅ Added {role} message")
    else:
        raise RuntimeError(f"{response.status_code}: {response.text}")


add_memory(
    "I need you to remember the key points of GDPR Article 15 - Right of Access...",
    "user",
)

add_memory(
    "I have noted the key provisions of GDPR Article 15. I will apply these when handling requests.",
    "assistant",
)
```

## Bulk Hydration Template

```python
hydration_batches = [
    {
        "context_key": "C45653",
        "messages": [
            {"role": "user", "message": "Customer John Doe, contract C45653, VIP status, prefers email."},
            {"role": "assistant", "message": "Logged: Customer John Doe (C45653) is VIP and prefers email."},
        ],
    },
    {
        "context_key": "C78901",
        "messages": [
            {"role": "user", "message": "Customer Jane Smith, contract C78901, standard tier, phone contact preferred."},
            {"role": "assistant", "message": "Logged: Customer Jane Smith (C78901) is standard tier and prefers phone calls."},
        ],
    },
]

for batch in hydration_batches:
    context_key = batch["context_key"]
    # Ensure the key exists first via POST /agent/.../context

    for entry in batch["messages"]:
        add_memory_to_context(context_key, entry["message"], entry["role"])
```

## Operational Tips

- **Create before hydrate:** the API rejects messages for missing keys with `400`
- **Preserve chronology:** inject messages in conversational order to avoid incoherent histories
- **Use roles intentionally:** keep `user` for inbound facts and `assistant` for stored knowledge; `system` is reserved for system-level instructions
- **Version knowledge:** append update messages rather than overwriting to maintain audit trails

For the human-friendly explanation of why hydration matters, see the [user documentation](../../users/memory/memory-hydration.md).
