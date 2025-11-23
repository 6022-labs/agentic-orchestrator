import { useAuth } from "@/contexts/auth-context";
import { initJwtAuth, finishJwtAuth } from "@/queries/api/auth.queries";
import {
  formatMessage,
  createSIWEConfig,
  type SIWEVerifyMessageArgs,
  type SIWECreateMessageArgs,
} from "@reown/appkit-siwe";

export const useSiweConfig = () => {
  const { signIn, signOut, user } = useAuth();

  return createSIWEConfig({
    getNonce: async (address?: string) => {
      if (!address) throw new Error("Address is required");

      const res = await initJwtAuth(address);
      const { nonce } = await res.data;
      return nonce;
    },
    getMessageParams: async () => ({
      domain: window.location.host,
      uri: window.location.origin,
      chains: [137],
      statement: "Please sign with your account",
    }),
    createMessage: ({ address, ...args }: SIWECreateMessageArgs) =>
      formatMessage(args, address),
    verifyMessage: async (args: SIWEVerifyMessageArgs) => {
      const res = await finishJwtAuth(args.message, args.signature);

      if (res.status >= 300) return false;
      const { jwt } = await res.data;
      signIn(jwt);
      return true;
    },
    getSession: async () => {
      if (user == null || user.exp == null || user.exp < Date.now() / 1000) {
        return {
          // To trigger a sign out
          address: "0x0000000000000000000000000000000000000000",
          chainId: -1,
        };
      }

      return user;
    },
    signOut: async () => {
      signOut();

      return true;
    },
  });
};
