# Contract fixtures

`artifacts/*.json` are trimmed hardhat artifacts (`{contractName, abi, bytecode}`)
that the integration tests deploy onto a throwaway anvil chain. They are
committed on purpose.

The abigen binding under `src/conversation_broker_blockchain/abi` cannot be used
instead — it carries the ABI only, no deployment bytecode.

## What is here

`AgentCollectionsManager`, `AgentCollectionsDescriptor`, `AgentCollectionV1`,
`AgentCollectionCreatorV1` and `AgentEnsRegistry` — the 6022 agentic contracts,
compiled from `agentic-smart-contracts` repository.

The broker only reads `AgentEnsRegistry.text()`, but the whole collection stack
has to exist to get a record in there: the registry derives an agent's subnode
from its collection, and only the agent's owner may publish its texts.

The same five artifacts are committed in
`agentic-agent-node/tests/agent_identity_blockchain_integration_tests/fixtures`;
refresh both together.

## Refreshing

Only needed when a 6022 agentic contract's ABI or behaviour changes in a way the
tests rely on.

Compile `agentic-smart-contracts` (`npx hardhat compile`), then keep only the
three fields — the full hardhat artifact carries megabytes of metadata:

```sh
python3 -c "
import json, sys
a = json.load(open(sys.argv[1]))
json.dump({'contractName': a['contractName'], 'abi': a['abi'], 'bytecode': a['bytecode']}, open(sys.argv[2], 'w'), indent=2)
" <hardhat>/artifacts/contracts/<Name>.sol/<Name>.json artifacts/<Name>.json
```
