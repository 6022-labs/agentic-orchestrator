# AgenticAI Orchestrator

This repository owns the code for the AgenticAI Orchestrator. The orchestrator mirrors on-chain agent identities (registry + chain synchronizer), brokers multi-agent conversations, and invokes each agent's own runtime (agentic-agent-node) discovered via its ENS url record in the AgentEnsRegistry and its `/.well-known/6022` card. Bridges (Slack, the web console) are thin translators that speak the broker's HTTP contract.

## How to run the orchestrator

You need Docker with Compose. Then run:

```bash
docker compose up
```

⚠ One of the containers can fail the first time you run it due to concurrent access to the database to create migrations. If this happens, just run the command again.

This runs the APIs on Polygon Amoy (public RPC), with Postgres, pgAdmin, Redis and a Grafana observability stack.

- Web console: `http://localhost:3000/`
- Gateway API: `http://localhost:5000/`
- pgAdmin: `http://localhost:5050/` (password `root`)
- Grafana: `http://localhost:3001/`

### On a local chain

```bash
docker compose -f docker-compose.yml -f docker-compose.local-chain.yml up
```

Adds an Anvil chain (id `31337`, RPC `http://localhost:8545`) with the 6022 contracts deployed from
[agentic-smart-contracts](https://github.com/6022-labs/agentic-smart-contracts); the synchronizer and broker switch to it.
Names on it don't resolve through public ENS. Reset the chain and the database together with `down -v` on the same two files.

## One-off cleanup after the agent-node revamp

Agent knowledge and compute moved into each agent's own node; their orphaned tables can be dropped with `scripts/teardown_part1.sql`.
