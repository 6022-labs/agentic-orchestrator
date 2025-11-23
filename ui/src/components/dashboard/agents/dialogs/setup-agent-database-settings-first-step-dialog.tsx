import { toast } from "sonner";
import { useAuth } from "@/contexts/auth-context";
import { useAgentDashboard } from "../contexts/agent-dashboard-context";
import { setupAgentDatabaseSettingsFirstStep } from "@/queries/api/agent.queries";
import SetupAgentDatabaseSettingsForm from "../forms/setup-agent-database-settings-first-step-form";
import { useSetupAgentDatabaseSettingsFirstStepDialogContext } from "../contexts/setup-agent-database-settings-first-step-dialog-context";
import { SetupAgentDatabaseFirstStepResult } from "../forms/setup-agent-database-settings-first-step-form/setup-agent-database-first-step-result-information";
import { SetupAgentDatabaseSettingsFirstStepFormSchemaType } from "../forms/setup-agent-database-settings-first-step-form/setup-agent-database-settings-first-step-form-schema";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";
import { extractValidationError } from "@/lib/api-errors";
import type { FormCompletionResult } from "@/models/api/validation-error";

export default function SetupAgentDatabaseSettingsFirstStepDialog() {
  const { jwt } = useAuth();
  const { fetchInternalAgents } = useAgentDashboard();
  const {
    agent,
    isOpen,
    setAgent,
    setIsOpen,
    hideCloseButton,
    isCloseButtonHidden,
  } = useSetupAgentDatabaseSettingsFirstStepDialogContext();

  const sendSetupDatabaseSettingsFirstStep = async (
    values: SetupAgentDatabaseSettingsFirstStepFormSchemaType
  ): Promise<
    FormCompletionResult<
      keyof SetupAgentDatabaseSettingsFirstStepFormSchemaType,
      SetupAgentDatabaseFirstStepResult
    >
  > => {
    try {
      const response = await setupAgentDatabaseSettingsFirstStep(
        jwt!,
        agent!.collectionAddress,
        agent!.collectionTokenId!,
        values
      );

      if (response.status >= 300) {
        throw new Error(
          "Failed to save database settings, status code: " + response.status
        );
      }

      toast.success("Database settings setting up successfully");

      return {
        result: response.data as SetupAgentDatabaseFirstStepResult,
      };
    } catch (error) {
      const validationError =
        extractValidationError<
          keyof SetupAgentDatabaseSettingsFirstStepFormSchemaType
        >(error);
      if (validationError) {
        return { error: validationError };
      }

      console.error("Error saving database settings:", error);
      toast.error("Failed to setup database settings");
      return {};
    }
  };

  const onFormFinished = () => {
    setIsOpen(false);
    setAgent(undefined);
    fetchInternalAgents();
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
      <DialogContent
        showCloseButton={!isCloseButtonHidden}
        onInteractOutside={(e) => {
          if (isCloseButtonHidden) {
            e.preventDefault();
          }
        }}
      >
        <DialogHeader>
          <DialogTitle>Setup agent database settings</DialogTitle>
          <DialogDescription>
            You don't known how to setup agent database? See our documentation.
          </DialogDescription>
        </DialogHeader>
        <SetupAgentDatabaseSettingsForm
          onFormFinished={onFormFinished}
          hideCloseButton={hideCloseButton}
          onFormComplete={sendSetupDatabaseSettingsFirstStep}
        />
      </DialogContent>
    </Dialog>
  );
}
