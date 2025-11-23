import { toast } from "sonner";
import { useAuth } from "@/contexts/auth-context";
import { extractAmazonPlatformRegionFromUrl } from "@/lib/utils";
import { setupAgentLLMSettings } from "@/queries/api/agent.queries";
import { useAgentDashboard } from "../contexts/agent-dashboard-context";
import SetupAgentLLMSettingsForm from "../forms/setup-agent-llm-settings-form";
import { LLMConnectionPlatformType } from "@/models/enums/llm-connection-platform-type";
import { useSetupAgentLLMSettingsDialogContext } from "../contexts/setup-agent-llm-settings-dialog-context";
import { SetupAgentLLMSettingsFormSchemaType } from "../forms/setup-agent-llm-settings-form/setup-agent-llm-settings-form-schema";
import {
  LLMConnectionPlatform,
  getLLMConnectionPlatformType,
} from "@/models/enums/llm-connection-platform";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";
import { extractValidationError } from "@/lib/api-errors";
import type { FormCompletionResult } from "@/models/api/validation-error";

export default function SetupAgentLLMSettingsDialog() {
  const { jwt } = useAuth();
  const {
    agent,
    isOpen,
    setAgent,
    setIsOpen,
    existingLLMSettings,
    isExistingLLMSettingsLoading,
  } = useSetupAgentLLMSettingsDialogContext();

  const { fetchInternalAgents } = useAgentDashboard();

  const sendSetupLLMSettings = async (
    values: SetupAgentLLMSettingsFormSchemaType
  ): Promise<
    FormCompletionResult<keyof SetupAgentLLMSettingsFormSchemaType>
  > => {
    try {
      const response = await setupAgentLLMSettings(
        jwt!,
        agent!.collectionAddress,
        agent!.collectionTokenId!,
        values
      );

      if (response.status >= 300) {
        throw new Error(
          "Failed to save llm settings, status code: " + response.status
        );
      }

      toast.success("LLM settings setting up successfully");
      fetchInternalAgents();
      setIsOpen(false);
      return {};
    } catch (error) {
      const validationError =
        extractValidationError<keyof SetupAgentLLMSettingsFormSchemaType>(
          error
        );
      if (validationError) {
        return { error: validationError };
      }

      console.error("Error saving llm settings:", error);
      toast.error("Failed to setup llm settings");
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
      <DialogContent loading={isExistingLLMSettingsLoading}>
        <DialogHeader>
          <DialogTitle>Setup agent LLM settings</DialogTitle>
          <DialogDescription>
            You don't known how to setup agent llm? See our documentation.
          </DialogDescription>
        </DialogHeader>
        <SetupAgentLLMSettingsForm
          defaultValues={
            existingLLMSettings
              ? {
                  model: existingLLMSettings.model,
                  platform: existingLLMSettings.platform,
                  platformType: getLLMConnectionPlatformType(
                    existingLLMSettings.platform
                  ),
                  apiKey: "",
                  baseUrl:
                    getLLMConnectionPlatformType(
                      existingLLMSettings.platform
                    ) === LLMConnectionPlatformType.SelfHost
                      ? existingLLMSettings.baseUrl
                      : undefined,
                  amazonPlatformRegion:
                    existingLLMSettings.platform ===
                      LLMConnectionPlatform.AmazonBedrock &&
                    existingLLMSettings.baseUrl
                      ? extractAmazonPlatformRegionFromUrl(
                          existingLLMSettings.baseUrl
                        )
                      : undefined,
                }
              : undefined
          }
          onFormComplete={sendSetupLLMSettings}
        />
      </DialogContent>
    </Dialog>
  );
}
