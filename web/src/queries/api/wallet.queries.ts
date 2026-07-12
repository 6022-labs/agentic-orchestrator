import axios from "axios"

// POST /wallet/verify-for-bridge — bind a bridge identity (e.g. Slack) to the
// signed-in wallet by signing the bridge-issued challenge.
export const verifyWalletForBridge = async (
  platform: string,
  challenge: string,
  signature: string
) => {
  return await axios.post("/wallet/verify-for-bridge", {
    platform,
    challenge,
    signature,
  })
}
