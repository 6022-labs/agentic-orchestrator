import { toast } from "sonner";
import { useAuth } from "@/contexts/auth-context";
import UpdateAgentPromptForm from "../forms/update-agent-prompt-form";
import { saveAgentPromptTemplate } from "@/queries/api/agent.queries";
import { UpdateAgentPromptFormSchemaType } from "../forms/update-agent-prompt-form/update-agent-prompt-form-schema";
import { useUpdateAgentPromptDialogContext } from "../contexts/update-agent-prompt-dialog-context";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";
import { extractValidationError } from "@/lib/api-errors";
import type { FormCompletionResult } from "@/models/api/validation-error";

export default function UpdateAgentPromptDialog() {
  const { jwt } = useAuth();
  const {
    agent,
    isOpen,
    setAgent,
    setIsOpen,
    isPromptTemplateLoading,
    agentPromptTemplateKeys,
    existingPromptsTemplates,
    reloadExistingPromptTemplate,
  } = useUpdateAgentPromptDialogContext();

  const sendSavePromptTemplate = async (
    values: UpdateAgentPromptFormSchemaType
  ): Promise<FormCompletionResult<keyof UpdateAgentPromptFormSchemaType>> => {
    try {
      const response = await saveAgentPromptTemplate(
        jwt!,
        agent!.collectionAddress,
        agent!.collectionTokenId!,
        values.key,
        values.content
      );

      if (response.status >= 300) {
        throw new Error(
          "Failed to save prompt template, status code: " + response.status
        );
      }

      toast.success("Prompt template updated successfully");

      await reloadExistingPromptTemplate(values.key);
      return {};
    } catch (error) {
      const validationError =
        extractValidationError<keyof UpdateAgentPromptFormSchemaType>(error);
      if (validationError) {
        return { error: validationError };
      }

      console.error("Error saving prompt template:", error);
      toast.error("Failed to update prompt template");
      return {};
    }
  };

  return (
    <Dialog
      open={isOpen}
      onOpenChange={(isOpen) => {
        setIsOpen(isOpen);

        if (!isOpen) {
          setAgent(undefined);
        }
      }}
    >
      <DialogContent loading={isPromptTemplateLoading}>
        <DialogHeader>
          <DialogTitle>Manage agent prompts</DialogTitle>
          <DialogDescription>
            You don't known how to manage agent prompts? See our documentation.
          </DialogDescription>
        </DialogHeader>
        <UpdateAgentPromptForm
          keys={agentPromptTemplateKeys}
          onFormComplete={sendSavePromptTemplate}
          defaultPromptsTemplates={existingPromptsTemplates}
        />
      </DialogContent>
    </Dialog>
  );
}
