# gateway_http_conversation_broker clients

`oapi_client.gen.go` is the generated HTTP transport the gateway uses to reach
the agent-swarm `conversation_broker_host`. It is produced by `oapi-codegen`
(pinned `v2.5.1`, `-generate types,client`) from
`agent-swarm/src/conversation_broker_host/docs/swagger.v3.yaml`.

Do not edit `oapi_client.gen.go` by hand. To refresh after the broker API changes:

```sh
sh agent-swarm/src/conversation_broker_host/docs/generate_docs.sh   # rebuild the broker OpenAPI spec
sh gateway/src/gateway_http_conversation_broker/clients/generate_client.sh
```

`../mappings` defines named types over the generated wire shapes and carries their
`ToUseCaseResponse()` mapping methods;
`../services/` wrap the generated `ClientWithResponses` behind the gateway port.
