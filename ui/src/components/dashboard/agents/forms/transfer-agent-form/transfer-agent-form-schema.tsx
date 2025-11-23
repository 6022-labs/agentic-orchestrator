import { ethers } from "ethers";
import { z } from "zod";

export type TransferAgentFormSchemaType = {
  address: string;
};

export const createTransferAgentFormSchema = () => {
  return z.object({
    address: z
      .string()
      .nonempty("Address is required")
      .refine((value) => {
        return ethers.isAddress(value) || value === "";
      }, "Invalid address format"),
  });
};
