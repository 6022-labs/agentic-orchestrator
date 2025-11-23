import { z } from "zod";
import { AgentContextMessageRole } from "@/models/enums/agent-context-message-role";

export type AddMessageToAgentContextFormSchemaType = {
  message: string;
  role: AgentContextMessageRole;
};

export const createAddMessageToAgentContextFormSchema = () => {
  return z.object({
    role: z.enum(AgentContextMessageRole, {
      error: () => ({
        message: "Role must be one of: user, system, assistant",
      }),
    }),
    message: z.string().min(2, {
      message: "Message must be at least 2 character long",
    }),
  });
};
