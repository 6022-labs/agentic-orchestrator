# conversation_broker_http_channels_and_integrations_bridge clients

`oapi_client.gen.go` is the generated HTTP transport the broker uses to post
completion callbacks to bridge internal hosts (today: `slack_bridge_internal_host`).
It is produced by `oapi-codegen` (pinned `v2.5.1`, `-generate types,client`) from
`channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/docs/swagger.v3.yaml`.

Each bridge's base URL lives in the `bridges.base_url` column, so
`client_factory.go` builds a `ClientWithResponses` per call-time base URL
instead of pinning one server at DI time.

Do not edit `oapi_client.gen.go` by hand. To refresh after the bridge API changes:

```sh
sh channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/docs/generate_docs.sh   # rebuild the bridge OpenAPI spec
sh agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/clients/generate_client.sh
```
