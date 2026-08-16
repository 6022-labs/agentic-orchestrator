# conversation_broker_http_agent_identity_registry clients

`oapi_client.gen.go` is the generated HTTP transport the conversation broker uses
to reach the agent-identity `registry_host`. It is produced by `oapi-codegen`
(pinned `v2.5.1`, `-generate types,client`) from
`agent-identity/src/registry_host/docs/swagger.v3.yaml`.

Do not edit `oapi_client.gen.go` by hand. To refresh after the registry API changes:

```sh
sh agent-identity/src/registry_host/docs/generate_docs.sh   # rebuild the registry OpenAPI spec
sh agent-swarm/src/conversation_broker_http_agent_identity_registry/clients/generate_client.sh
```

`../mappings` defines named types over the generated wire shapes and carries their
mapping methods; `../services/` wrap the generated `ClientWithResponses` behind the
broker port.
