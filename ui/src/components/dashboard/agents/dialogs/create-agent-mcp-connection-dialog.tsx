import { toast } from "sonner";
import { useAuth } from "@/contexts/auth-context";
import { createAgentMCPConnection } from "@/queries/api/agent.queries";
import AddAgentMCPConnectionForm from "../forms/create-agent-mcp-connection-form";
import { useAgentMCPConnectionsDialogContext } from "../contexts/agent-mcp-connections-dialog-context";
import { useCreateAgentMCPConnectionDialogContext } from "../contexts/create-agent-mcp-connection-dialog-context";
import { CreateAgentMCPConnectionFormSchemaType } from "../forms/create-agent-mcp-connection-form/create-agent-mcp-connection-form-schema";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { extractValidationError } from "@/lib/api-errors";
import type { FormCompletionResult } from "@/models/api/validation-error";

export default function CreateAgentMCPConnectionDialog() {
  const { jwt } = useAuth();
  const {
    agent,
    fetchMCPConnections,
    setIsOpen: setMCPSettingsDialogOpen,
  } = useAgentMCPConnectionsDialogContext();
  const { isOpen, setIsOpen } = useCreateAgentMCPConnectionDialogContext();

  const sendCreate = async (
    values: CreateAgentMCPConnectionFormSchemaType
  ): Promise<
    FormCompletionResult<keyof CreateAgentMCPConnectionFormSchemaType>
  > => {
    try {
      const response = await createAgentMCPConnection(
        jwt!,
        agent!.collectionAddress,
        agent!.collectionTokenId!,
        values
      );

      if (response.status >= 300) {
        throw new Error(
          "Failed to create MCP connection, status code: " + response.status
        );
      }

      toast.success("MCP connection created successfully");
      await fetchMCPConnections();

      setIsOpen(false);
      setMCPSettingsDialogOpen(true);
      return {};
    } catch (error) {
      const validationError =
        extractValidationError<keyof CreateAgentMCPConnectionFormSchemaType>(
          error
        );
      if (validationError) {
        return { error: validationError };
      }

      console.error("Failed to create MCP connection:", error);
      toast.error("Failed to create MCP connection");
      return {};
    }
  };

  return (
    <Dialog
      open={isOpen}
      onOpenChange={(open) => {
        if (!open) {
          setMCPSettingsDialogOpen(true);
        }
        setIsOpen(open);
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add MCP</DialogTitle>
          <DialogDescription>
            <span>Configure a new MCP connection for agent </span>
            <span className="font-bold">{agent?.name}</span>
          </DialogDescription>
        </DialogHeader>
        <AddAgentMCPConnectionForm onFormComplete={sendCreate} />
      </DialogContent>
    </Dialog>
  );
}
