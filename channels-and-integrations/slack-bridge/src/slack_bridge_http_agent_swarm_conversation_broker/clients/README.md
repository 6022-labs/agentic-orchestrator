# slack_bridge_http_agent_swarm_conversation_broker clients

`oapi_client.gen.go` is the generated HTTP transport the slack bridge uses to
reach the agent-swarm `conversation_broker_host`. It is produced by
`oapi-codegen` (pinned `v2.5.1`, `-generate types,client`) from
`agent-swarm/src/conversation_broker_host/docs/swagger.v3.yaml`.

Do not edit `oapi_client.gen.go` by hand. To refresh after the broker API changes:

```sh
sh agent-swarm/src/conversation_broker_host/docs/generate_docs.sh   # rebuild the broker OpenAPI spec
sh channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/clients/generate_client.sh
```

`../helpers` holds the shared pointer/uuid conversions; `../services/` wrap the
generated `ClientWithResponses` directly behind the bridge ports (the services
only consume ids, so there is no `mappings/` layer here).
