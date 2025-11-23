import { toast } from "sonner";
import { useAuth } from "@/contexts/auth-context";
import { addMessageToContext } from "@/queries/api/agent.queries";
import AddMessageToAgentContextForm from "../forms/add-message-to-agent-context-form";
import { useAgentContextsDialogContext } from "../contexts/agent-contexts-dialog-context";
import { useAddMessageToAgentContextDialogContext } from "../contexts/add-message-to-agent-context-dialog-context";
import { AddMessageToAgentContextFormSchemaType } from "../forms/add-message-to-agent-context-form/add-message-to-agent-context-form-schema";
import { extractValidationError } from "@/lib/api-errors";
import type { FormCompletionResult } from "@/models/api/validation-error";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";

export default function AddMessageToAgentContextDialog() {
  const { jwt } = useAuth();
  const { isOpen, setIsOpen, context } =
    useAddMessageToAgentContextDialogContext();
  const { agent, setIsOpen: setAgentContextsDialogOpen } =
    useAgentContextsDialogContext();

  const sendAddMessageToContext = async (
    values: AddMessageToAgentContextFormSchemaType
  ): Promise<
    FormCompletionResult<keyof AddMessageToAgentContextFormSchemaType>
  > => {
    try {
      const response = await addMessageToContext(
        jwt!,
        agent!.collectionAddress,
        agent!.collectionTokenId!,
        context,
        values.message,
        values.role
      );

      if (response.status >= 300) {
        throw new Error(
          "Failed to add message to agent context: " + response.statusText
        );
      }

      toast.success("Agent context message added successfully");
      return {};
    } catch (error) {
      const validationError =
        extractValidationError<keyof AddMessageToAgentContextFormSchemaType>(
          error
        );
      if (validationError) {
        return { error: validationError };
      }

      console.error("Failed to add message to agent context:", error);
      toast.error("Failed to add message to agent context");
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
          <DialogTitle>Add message to agent context</DialogTitle>
          <DialogDescription>
            <span>Add a message to the context </span>
            <span className="font-bold">{context}</span>
            <span> for agent </span>
            <span className="font-bold">{agent?.name}</span>.
          </DialogDescription>
        </DialogHeader>
        <AddMessageToAgentContextForm
          onFormComplete={sendAddMessageToContext}
        />
      </DialogContent>
    </Dialog>
  );
}
