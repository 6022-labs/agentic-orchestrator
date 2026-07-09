# AgenticAI Orchestrator

This repository owns the code for the AgenticAI Orchestrator. The orchestrator mirrors on-chain agent identities (registry + chain synchronizer), brokers multi-agent conversations, and invokes each agent's own runtime (agentic-agent-node) discovered via ENS CCIP-Read and its `/.well-known/6022` card. Bridges (Slack, the web console) are thin translators that speak the broker's HTTP contract.

## How to run the orchestrator

You will need to have `docker` and `docker-compose` installed in your machine. Then, you can run the following command:

```bash
docker-compose up
```

⚠ One of the containers can fail the first time you run it due to concurrent access to the database to create migrations. If this happens, just run the command again.

This command will mount a database, a pgadmin (to visualize the database), a hashicorp vault, a local blockchain and the APIs.

- You can access the web console at `http://localhost:3000/`
- You can access the gateway API at `http://localhost:5000/`
- You can access the blockchain at `http://localhost:8545/`
- You can access the vault at `http://localhost:8200/` (use `wBGKZCJjfIGqIKB3HMmvcQVt` as passwords)
- You can access the pgadmin at `http://localhost:5050/` (use `root` as passwords)

## One-off cleanup after the agent-node revamp

Agent knowledge and compute moved into each agent's own node; their orphaned tables can be dropped with `scripts/teardown_part1.sql`.
