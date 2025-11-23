import axios from "axios";

export const verifyWalletForBridge = async (
  platform: string,
  challenge: string,
  signature: string,
  jwt: string
) => {
  return await axios.post(
    "/wallet/verify-for-bridge",
    {
      platform,
      challenge,
      signature,
    },
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    }
  );
};
