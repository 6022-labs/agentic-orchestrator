import { z } from "zod";

export type CreateAgentContextFormSchemaType = {
  key: string;
};

export const createCreateAgentContextFormSchema = () => {
  return z.object({
    key: z
      .string()
      .min(2, {
        message: "Key must be at least 2 characters long",
      })
      .refine(
        (data) => {
          return !data.startsWith("conversation-context-discovery");
        },
        {
          message: "Key cannot start with 'conversation-context-discovery'",
        }
      ),
  });
};
