import { toast } from "sonner";
import { useAuth } from "@/contexts/auth-context";
import { createAgentContext } from "@/queries/api/agent.queries";
import CreateAgentContextForm from "../forms/create-agent-context-form";
import { useCreateAgentContextDialogContext } from "../contexts/create-agent-context-dialog-context";
import { CreateAgentContextFormSchemaType } from "../forms/create-agent-context-form/create-agent-context-form-schema";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";
import { useAgentContextsDialogContext } from "../contexts/agent-contexts-dialog-context";
import { extractValidationError } from "@/lib/api-errors";
import type { FormCompletionResult } from "@/models/api/validation-error";

export default function CreateAgentContextDialog() {
  const { jwt } = useAuth();
  const { isOpen, setIsOpen } = useCreateAgentContextDialogContext();
  const { agent, setIsOpen: setAgentContextsDialogOpen } =
    useAgentContextsDialogContext();

  const sendCreateContext = async (
    values: CreateAgentContextFormSchemaType
  ): Promise<FormCompletionResult<keyof CreateAgentContextFormSchemaType>> => {
    try {
      const response = await createAgentContext(
        jwt!,
        agent!.collectionAddress,
        agent!.collectionTokenId!,
        values.key
      );

      if (response.status >= 300) {
        throw new Error(
          "Failed to create context, status code: " + response.status
        );
      }

      toast.success("Agent context created successfully");
      setIsOpen(false);
      return {};
    } catch (error) {
      const validationError =
        extractValidationError<keyof CreateAgentContextFormSchemaType>(error);
      if (validationError) {
        return { error: validationError };
      }

      console.error("Failed to create agent context:", error);
      toast.error("Failed to create agent context");
      return {};
    }
  };

  return (
    <Dialog
      open={isOpen}
      onOpenChange={(isOpen) => {
        setIsOpen(isOpen);

        if (!isOpen) {
          setAgentContextsDialogOpen(true);
        }
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogHeader>
            <DialogTitle>Create agent context</DialogTitle>
            <DialogDescription>
              <span>Create a context for agent </span>
              <span className="font-bold">{agent?.name}</span>.
            </DialogDescription>
          </DialogHeader>
        </DialogHeader>
        <CreateAgentContextForm onFormComplete={sendCreateContext} />
      </DialogContent>
    </Dialog>
  );
}
