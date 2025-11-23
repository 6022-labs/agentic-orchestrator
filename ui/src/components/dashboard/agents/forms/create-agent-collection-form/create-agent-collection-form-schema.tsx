import { z } from "zod";
import { Eip1193Provider } from "ethers";
import {
  collectionIdByName,
  collectionIdBySymbol,
} from "@/queries/contracts/agent-collections-manager";

export type CreateAgentCollectionFormSchemaType = {
  name: string;
  symbol: string;
  description: string;
};

export const createCreateAgentCollectionFormSchema = (
  walletProvider: Eip1193Provider
) => {
  return z.object({
    name: z
      .string()
      .min(1, {
        message: "Name must be at least 1 character.",
      })
      .superRefine(async (name, ctx) => {
        try {
          const collectionId = await collectionIdByName(name, walletProvider);

          if (collectionId != 0n) {
            ctx.addIssue({
              code: "custom",
              message: "Name is already taken.",
            });
          }
        } catch (e) {
          console.error(e);
          ctx.addIssue({
            code: "custom",
            message: "Cannot check if name is taken.",
          });
        }
      }),
    symbol: z
      .string()
      .min(1, {
        message: "Symbol must be at least 1 character.",
      })
      .superRefine(async (name, ctx) => {
        try {
          const collectionId = await collectionIdBySymbol(name, walletProvider);

          if (collectionId != 0n) {
            ctx.addIssue({
              code: "custom",
              message: "Symbol is already taken.",
            });
          }
        } catch (e) {
          console.error(e);
          ctx.addIssue({
            code: "custom",
            message: "Cannot check if symbol is taken.",
          });
        }
      }),
    description: z.string().min(1, {
      message: "Description must be at least 1 character.",
    }),
  });
};
