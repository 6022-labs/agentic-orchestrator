// Centralised access to the Vite env with Amoy-friendly fallbacks, so the app
// runs out of the box and a .env only needs to override what differs.

const env = import.meta.env

export const config = {
  appKitProjectId:
    env.VITE_APP_KIT_PROJECT_ID || "6e030f5fb755b144ec2cb0d10f2572ae",
  // Relative by default: nginx (container) and the Vite dev proxy both route
  // /api to the gateway, so the same bundle works behind any host.
  gatewayApiBaseUrl: env.VITE_GATEWAY_API_BASE_URL || "/api",
  // Multiplies the provider's suggested gas fees for on-chain writes (>= 1).
  txGasMultiplier: Math.max(Number(env.VITE_TX_GAS_MULTIPLIER) || 1, 1),
}
