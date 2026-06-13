import { toast } from "sonner";
import { useWallet } from "@/contexts/wallet-provider-context";
import { useAgentDashboard } from "../contexts/agent-dashboard-context";
import UpsertAgentAttributeForm from "../forms/upsert-agent-attribute-form";
import { addOrUpdateAgentAttribute } from "@/queries/contracts/agent-collection";
import { useUpsertAgentAttributeDialogContext } from "../contexts/upsert-agent-attribute-dialog-context";
import { UpsertAgentAttributeFormSchema } from "../forms/upsert-agent-attribute-form/upsert-agent-attribute-form-schema";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";
import type { FormCompletionResult } from "@/models/api/validation-error";

export default function UpsertAgentAttributeDialog() {
  const { agent, isOpen, setAgent, setIsOpen } =
    useUpsertAgentAttributeDialogContext();
  const { fetchInternalAgents, fetchExternalAgents } = useAgentDashboard();

  const { wallet } = useWallet();

  const upsertAgentAttribute = async (
    values: UpsertAgentAttributeFormSchema
  ): Promise<FormCompletionResult<keyof UpsertAgentAttributeFormSchema>> => {
    try {
      await addOrUpdateAgentAttribute(
        agent!.collectionAddress,
        agent!.collectionTokenId!,
        values.key,
        values.value,
        wallet!
      );

      toast.success("Agent attribute upserted successfully");
      fetchInternalAgents();
      fetchExternalAgents();
      setIsOpen(false);
      return {};
    } catch (error) {
      console.error("Error upserting agent attribute:", error);
      toast.error("Failed to upsert agent attribute");
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
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Upsert agent attribute</DialogTitle>
          <DialogDescription>
            You don't known how to upsert agents attribute? See our
            documentation.
          </DialogDescription>
        </DialogHeader>
        <UpsertAgentAttributeForm
          onFormComplete={upsertAgentAttribute}
          existingAttributes={agent?.attributes || {}}
        />
      </DialogContent>
    </Dialog>
  );
}
