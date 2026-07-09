import { toast } from "sonner";
import TransferAgentForm from "../forms/transfer-agent-form";
import { useWallet } from "@/contexts/wallet-provider-context";
import { transferFrom } from "@/queries/contracts/agent-collection";
import { useAgentDashboard } from "../contexts/agent-dashboard-context";
import { useTransferAgentDialogContext } from "../contexts/transfer-agent-dialog-context";
import { TransferAgentFormSchemaType } from "../forms/transfer-agent-form/transfer-agent-form-schema";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";
import type { FormCompletionResult } from "@/models/api/validation-error";

export default function TransferAgentDialog() {
  const { agent, isOpen, setAgent, setIsOpen } =
    useTransferAgentDialogContext();
  const { fetchExternalAgents } = useAgentDashboard();

  const { wallet } = useWallet();

  const transferAgent = async (
    values: TransferAgentFormSchemaType
  ): Promise<FormCompletionResult<keyof TransferAgentFormSchemaType>> => {
    try {
      await transferFrom(
        agent!.collectionAddress,
        agent!.collectionTokenId!,
        values.address,
        wallet!
      );

      toast.success("Agent transferred successfully");
      fetchExternalAgents();
      return {};
    } catch (error) {
      console.error("Error transferring agent:", error);
      toast.error("Failed to transfer agent");
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
          <DialogTitle>Transfer agent</DialogTitle>
          <DialogDescription>
            You don't known how to transfer agents? See our documentation.
          </DialogDescription>
        </DialogHeader>
        <TransferAgentForm onFormComplete={transferAgent} />
      </DialogContent>
    </Dialog>
  );
}
