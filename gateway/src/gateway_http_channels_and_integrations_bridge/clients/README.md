# gateway_http_channels_and_integrations_bridge clients

`oapi_client.gen.go` is the generated HTTP transport the gateway uses to reach
platform bridge internal hosts (today: `slack_bridge_internal_host`). It is
produced by `oapi-codegen` (pinned `v2.5.1`, `-generate types,client`) from
`channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/docs/swagger.v3.yaml`.

Do not edit `oapi_client.gen.go` by hand. To refresh after the bridge API changes:

```sh
sh channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/docs/generate_docs.sh   # rebuild the bridge OpenAPI spec
sh gateway/src/gateway_http_channels_and_integrations_bridge/clients/generate_client.sh
```

`../services/` wrap the generated `ClientWithResponses` behind the gateway port,
one client per configured bridge platform.
