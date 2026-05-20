import { useMemo } from "react";
import { useAuth } from "@/contexts/auth-context";
import { initJwtAuth, finishJwtAuth } from "@/queries/api/auth.queries";
import {
  formatMessage,
  createSIWEConfig,
  type SIWEVerifyMessageArgs,
  type SIWECreateMessageArgs,
} from "@reown/appkit-siwe";

export const useSiweConfig = () => {
  const { signIn, signOut } = useAuth();

  return useMemo(
    () =>
      createSIWEConfig({
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
          const token = localStorage.getItem("jwt");
          if (!token) return null;

          try {
            const payload = JSON.parse(atob(token.split(".")[1]));
            if (!payload.exp || payload.exp < Date.now() / 1000) return null;

            return { address: payload.sub, chainId: Number(payload.cId) };
          } catch {
            return null;
          }
        },
        signOut: async () => {
          signOut();

          return true;
        },
      }),
    [signIn, signOut],
  );
};
