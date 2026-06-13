import { toast } from "sonner";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { extractValidationError } from "@/lib/api-errors";
import { createAgentFirstStep } from "@/queries/api/agent.queries";
import { useAgentDashboard } from "../contexts/agent-dashboard-context";
import { getAllAgentRolesExceptClone } from "@/models/enums/agent-role";
import type { FormCompletionResult } from "@/models/api/validation-error";
import CreateAgentFirstStepForm from "../forms/create-agent-first-step-form";
import { useCreateAgentFirstStepDialog } from "../contexts/create-agent-first-step-dialog-context";
import { CreateAgentFirstStepFormSchemaType } from "../forms/create-agent-first-step-form/create-agent-first-step-form-schema";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogTrigger,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";

export default function CreateAgentFirstStepDialog() {
  const {
    isOpen,
    setIsOpen,
    defaultValues,
    disabledValues,
    accessibleRoles,
    setDefaultValues,
    setDisabledValues,
    setAccessibleRoles,
  } = useCreateAgentFirstStepDialog();

  const { fetchInternalAgents } = useAgentDashboard();

  const onFormComplete = async (
    values: CreateAgentFirstStepFormSchemaType
  ): Promise<
    FormCompletionResult<keyof CreateAgentFirstStepFormSchemaType>
  > => {
    try {
      const response = await createAgentFirstStep(values);
      if (response.status !== 200) {
        throw new Error("Received non-200 response");
      }
      toast.success("Agent information registered successfully");
      fetchInternalAgents();
      setIsOpen(false);
      return {};
    } catch (error) {
      const validationError =
        extractValidationError<keyof CreateAgentFirstStepFormSchemaType>(error);
      if (validationError) {
        return { error: validationError };
      }

      console.error("Error registering agent information:", error);
      toast.error("Failed to register agent information");
      return {};
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button
          type="button"
          onClick={() => {
            setDefaultValues({});
            setDisabledValues({});
            setAccessibleRoles(getAllAgentRolesExceptClone());
          }}
        >
          <Plus className="mr-1" />
          Create Agent
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create an agent</DialogTitle>
          <DialogDescription>
            You don't known how to create an agent? See our documentation.
          </DialogDescription>
        </DialogHeader>
        <CreateAgentFirstStepForm
          defaultValues={defaultValues}
          disabledValues={disabledValues}
          onFormComplete={onFormComplete}
          accessibleRoles={accessibleRoles}
        />
      </DialogContent>
    </Dialog>
  );
}
