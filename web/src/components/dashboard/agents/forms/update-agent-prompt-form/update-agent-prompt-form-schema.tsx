import { z } from "zod";
import { AgentPromptTemplateKey } from "@/models/enums/agent-prompt-template-key";

export type UpdateAgentPromptFormSchemaType = {
  key: AgentPromptTemplateKey;
  content: string;
};

export const createUpdateAgentPromptFormSchema = () => {
  return z.object({
    key: z.enum(AgentPromptTemplateKey),
    content: z.string(),
  });
};
