# Contract fixtures

`artifacts/*.json` are trimmed hardhat artifacts (`{contractName, abi, bytecode}`)
that the synchronizer integration tests deploy onto a throwaway anvil chain. They
are committed on purpose.

The abigen bindings under `src/registry_common/abi` cannot be used instead — they
carry the ABI only, no deployment bytecode.

## What is here

`AgentCollectionsManager`, `AgentCollectionsDescriptor`, `AgentCollectionV1` and
`AgentCollectionCreatorV1` — the 6022 agentic contracts, compiled from
`agentic-smart-contracts` repository.

## Refreshing

Only needed when a contract's ABI or behaviour changes in a way the tests rely on.

Compile `agentic-smart-contracts` (`npx hardhat compile`), then keep only the
three fields — the full hardhat artifact carries megabytes of metadata:

```sh
python3 -c "
import json, sys
a = json.load(open(sys.argv[1]))
json.dump({'contractName': a['contractName'], 'abi': a['abi'], 'bytecode': a['bytecode']}, open(sys.argv[2], 'w'), indent=2)
" <hardhat>/artifacts/contracts/<Name>.sol/<Name>.json artifacts/<Name>.json
```
