---
sidebar_position: 2
title: Authentication
---

# Authentication

Every write request to the orchestrator API is gated by a JWT that you obtain by signing a Sign-In With Ethereum (SIWE) message with the wallet that owns the agent. There is no username/password — your wallet **is** the credential.

## How it works

```
┌─────────┐                   ┌──────────────────┐
│ Client  │                   │   Orchestrator   │
│ (wallet)│                   │     gateway      │
└────┬────┘                   └────────┬─────────┘
     │                                 │
     │ 1. POST /auth/jwt/init          │
     │    { walletAddress, chainId }   │
     │────────────────────────────────►│
     │                                 │
     │           nonce + SIWE message  │
     │◄────────────────────────────────│
     │                                 │
     │ 2. Sign the SIWE message        │
     │    locally with your wallet     │
     │                                 │
     │ 3. POST /auth/jwt/finish        │
     │    { walletAddress, signature } │
     │────────────────────────────────►│
     │                                 │
     │            JWT (Bearer token)   │
     │◄────────────────────────────────│
     │                                 │
     │ 4. Use Authorization: Bearer    │
     │    on every subsequent call     │
     │────────────────────────────────►│
```

The orchestrator never sees your private key. It only verifies the ECDSA signature on a nonce that it itself issued, then mints a JWT bound to the wallet address and chain ID.

## Endpoints

### `POST /auth/jwt/init`

Start a session. The server returns a fresh nonce and a SIWE message you must sign locally.

**Body**

```json
{
  "walletAddress": "0xABC...",
  "chainId": 137
}
```

`chainId` is `137` for Polygon mainnet, `80002` for Amoy testnet.

**Response**

```json
{
  "nonce": "abcdef0123...",
  "message": "<your domain> wants you to sign in with your Ethereum account: 0xABC...\n\n..."
}
```

### `POST /auth/jwt/finish`

Submit the signed message and receive a JWT.

**Body**

```json
{
  "walletAddress": "0xABC...",
  "chainId": 137,
  "signature": "0x<65-byte hex>"
}
```

**Response**

```json
{
  "token": "eyJhbGciOi..."
}
```

The server verifies that the signature matches the wallet address and the original nonce. If anything is off (wrong signer, expired nonce, malformed message), the request is rejected with `401`.

## Using the token

Send the token in the `Authorization` header on every subsequent call:

```bash
curl -H "Authorization: Bearer eyJhbGciOi..." \
     https://<your-orchestrator-api-host>/agent/0xCOLLECTION/1/contexts/keys
```

## JWT contents

The token is a standard JWT signed by the orchestrator with the secret configured as `JWT__HASH`. Its claims include:

| Claim | Meaning |
|---|---|
| `sub` | Wallet address that signed the SIWE message |
| `cId` | Chain ID the token is bound to (`137` or `80002`) |
| `iat` | Issued-at timestamp |
| `nbf` | Not-before timestamp |
| `exp` | Expiry timestamp — controlled by `JWT__LIVING_TIME` (seconds) |

A token issued for chain `137` is **not** accepted by an orchestrator instance running against Amoy, and vice versa.

## Generating a JWT programmatically

You have two options.

### Option A — call the API directly

Run the two-step `/auth/jwt/init` → `/auth/jwt/finish` flow yourself. Any web3 library that can sign a SIWE message works (`ethers`, `viem`, `web3.py`, etc.).

```python
from eth_account import Account
from eth_account.messages import encode_defunct
import requests

WALLET = Account.from_key("0x<your private key>")
BASE_URL = "https://<your-orchestrator-api-host>"
CHAIN_ID = 137

# 1. Get nonce + SIWE message
init = requests.post(f"{BASE_URL}/auth/jwt/init", json={
    "walletAddress": WALLET.address,
    "chainId": CHAIN_ID,
}).json()

# 2. Sign the SIWE message locally
message = init["message"]
signed = WALLET.sign_message(encode_defunct(text=message))

# 3. Exchange the signature for a JWT
finish = requests.post(f"{BASE_URL}/auth/jwt/finish", json={
    "walletAddress": WALLET.address,
    "chainId": CHAIN_ID,
    "signature": signed.signature.hex(),
}).json()

token = finish["token"]
```

### Option B — use the helper script

The 6022 organisation publishes a small Python tool that wraps the same flow:

```bash
git clone https://github.com/6022protocol/py-generate-jwt
cd py-generate-jwt

python generate_jwt.py \
  --private-key 0x<your private key> \
  --chain-id 137 \
  --api-base-url https://<your-orchestrator-api-host>/
```

| Parameter | Description |
|---|---|
| `--private-key` | Wallet private key (owner of the agent) |
| `--chain-id` | `137` for Polygon mainnet, `80002` for Amoy testnet |
| `--api-base-url` | Orchestrator API base URL, with trailing slash |

The script prints a Bearer token you can paste into an `Authorization` header.

:::warning
Use a dedicated wallet for API access. Never check a real owner key into a script, a CI variable, or a `.env` file you don't control.
:::

## Authorization model

A valid JWT proves that you control a wallet — it does **not** prove you own a particular agent. Per-request authorization is enforced by the orchestrator: when you call any `/agent/:agentCollectionAddress/:agentCollectionTokenId/...` endpoint, the gateway checks that the wallet in your token's `sub` claim is the on-chain owner of that NFT. If you transfer the agent to a different wallet, your existing JWT immediately stops working for that agent.

## Token lifetime

JWTs expire after `JWT__LIVING_TIME` seconds (configured on the orchestrator). When a token expires, repeat the init/finish flow to get a new one. There is no refresh-token endpoint — the SIWE flow is cheap enough that re-issuing is the supported path.

## See also

- [API Reference overview](./intro.md) — every endpoint requires the `Authorization: Bearer` header described here.
- [Context Primary Keys API](./context-primary-keys.md) — the main resource you'll call once you have a token.
- [Memory Hydration API](./memory-hydration.md) — bulk knowledge injection.
- [User-facing orchestrator intro](../../users/orchestrator/intro.md) — for the one-paragraph "your wallet is your login" explanation.
