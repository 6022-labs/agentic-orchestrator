# Contract fixtures

`artifacts/*.json` are trimmed hardhat artifacts (`{contractName, abi, bytecode}`)
that the ENS integration tests deploy onto a throwaway anvil chain. They are
committed on purpose.

The abigen bindings under `src/conversation_broker_blockchain/abi` cannot be used
instead — they carry the ABI only, no deployment bytecode.

## What is here

`AgentEnsOffchainResolver` — the 6022 agentic CCIP-Read resolver, compiled from
`agentic-smart-contracts` repository.

`EnsRegistry` — the canonical ENS registry, taken from the
`@ensdomains/ens-contracts` npm artifacts. This is ENS's own contract, not
6022's `AgentEnsRegistry`: the tests deploy the real registry so the resolution
walk exercises real node ownership and `setSubnodeOwner` semantics rather than a
stand-in.

## Refreshing

Rarely needed: ENS registry behaviour is fixed by the standard, and the resolver
only changes when the 6022 agentic CCIP-Read contract does.

Compile the source, then keep only the three fields — the full hardhat artifact
carries megabytes of metadata:

```sh
python3 -c "
import json, sys
a = json.load(open(sys.argv[1]))
json.dump({'contractName': a['contractName'], 'abi': a['abi'], 'bytecode': a['bytecode']}, open(sys.argv[2], 'w'), indent=2)
" <source-artifact>.json artifacts/<Name>.json
```
