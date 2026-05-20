import axios from "axios";

export const verifyWalletForBridge = async (
  platform: string,
  challenge: string,
  signature: string
) => {
  return await axios.post("/wallet/verify-for-bridge", {
    platform,
    challenge,
    signature,
  });
};
