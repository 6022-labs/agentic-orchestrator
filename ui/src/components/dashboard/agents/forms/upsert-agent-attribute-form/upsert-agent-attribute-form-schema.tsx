import { z } from "zod";

export type UpsertAgentAttributeFormSchema = {
  key: string;
  value: string;
};

export const createUpsertAgentAttributeFormSchema = () => {
  return z.object({
    key: z.string().min(1, "Key is required"),
    value: z.string().min(1, "Value is required"),
  });
};
