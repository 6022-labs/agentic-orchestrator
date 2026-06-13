import { z } from "zod";
import { Eip1193Provider, ethers } from "ethers";
import { AgentRole } from "@/models/enums/agent-role";
import {
  nameToTokenId,
  walletToTokenId,
} from "@/queries/contracts/agent-collection";

export type CreateAgentFirstStepDisabledValues = {
  collectionAddress?: boolean;
  name?: boolean;
  image?: boolean;
  role?: boolean;
  mission?: boolean;
  walletAddress?: boolean;
  cloneOfAgentId?: boolean;
};

export type CreateAgentFirstStepFormSchemaType = {
  collectionAddress: string;
  name: string;
  image: File;
  role: AgentRole;
  mission: string;
  walletAddress?: string;
  cloneOfAgentId?: number;
};

const validImageTypes = ["image/png", "image/gif", "image/jpg", "image/jpeg"];

export const createCreateAgentFirstStepFormSchema = (
  eip1193Provider: Eip1193Provider
) => {
  const baseSchema = z.object({
    collectionAddress: z
      .string()
      .min(42, {
        message: "Collection address must be a valid Ethereum address.",
      })
      .max(42, {
        message: "Collection address must be a valid Ethereum address.",
      })
      .refine((value) => /^0x[a-fA-F0-9]{40}$/.test(value), {
        message: "Collection address must be a valid Ethereum address.",
      }),
    name: z.string().min(1, {
      message: "Name must be at least 1 character.",
    }),
    image: z
      .instanceof(File)
      .refine((file) => validImageTypes.includes(file.type), {
        message: `Invalid image type. Allowed types are: ${validImageTypes.join(
          ", "
        )}.`,
      })
      .refine((file) => file.size <= 5 * 1024 * 1024, {
        message: "Image size must be less than 5MB.",
      }),
    role: z.enum(AgentRole).superRefine((value, ctx) => {
      if (value === AgentRole.None) {
        ctx.addIssue({
          code: "custom",
          message: "Agent role is required.",
        });
      }
    }),
    mission: z.string().min(1, {
      message: "Mission must be at least 1 character.",
    }),
    walletAddress: z.string().optional(),
    cloneOfAgentId: z.number().optional(),
  });

  return baseSchema.superRefine(async (data, ctx) => {
    // Only validate Human wallet address when role is explicitly provided as Human
    if (data.role === AgentRole.Human) {
      if (!data.walletAddress) {
        ctx.addIssue({
          path: ["walletAddress"],
          code: "custom",
          message: "Wallet address is required for Human agents.",
        });
      } else if (!ethers.isAddress(data.walletAddress)) {
        ctx.addIssue({
          path: ["walletAddress"],
          code: "custom",
          message: "Wallet address must be a valid Ethereum address.",
        });
      } else {
        const tokenId = await walletToTokenId(
          data.collectionAddress,
          data.walletAddress,
          eip1193Provider
        );

        if (tokenId != 0n) {
          ctx.addIssue({
            path: ["walletAddress"],
            code: "custom",
            message: "Wallet address is already taken.",
          });
        }
      }
    }

    if (data.role === AgentRole.Clone && !data.cloneOfAgentId) {
      ctx.addIssue({
        path: ["cloneOfAgentId"],
        code: "custom",
        message: "Clone of Agent ID is required for Clone agents.",
      });
    }

    // If either collectionAddress or name is missing (due to partial), skip uniqueness check
    if (!data.collectionAddress || !data.name) return;

    try {
      const tokenId = await nameToTokenId(
        data.collectionAddress,
        data.name,
        eip1193Provider
      );

      if (tokenId != 0n) {
        ctx.addIssue({
          path: ["name"],
          code: "custom",
          message: "Name is already taken.",
        });
      }
    } catch (error) {
      console.error("Error checking if name is taken:", error);
      ctx.addIssue({
        code: "custom",
        message: "Failed to check if name is taken.",
      });
    }
  });
};
