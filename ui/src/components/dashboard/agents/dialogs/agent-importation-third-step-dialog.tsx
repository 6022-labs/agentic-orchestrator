import React from "react";
import { toast } from "sonner";
import { strToBytes } from "@/lib/utils";
import { Loader2Icon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Eip1193Provider, ethers } from "ethers";
import { useAuth } from "@/contexts/auth-context";
import { useAppKitProvider } from "@reown/appkit/react";
import { useAgentDashboard } from "../contexts/agent-dashboard-context";
import { updateAgentSigner } from "@/queries/contracts/agent-collection";
import { importAgentThirdStep } from "@/queries/api/import-agent.queries";
import { useAgentImportationThirdStepDialog } from "../contexts/agent-importation-third-step-dialog-context";
import {
  Dialog,
  DialogTitle,
  DialogHeader,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";

export default function AgentImportationThirdStepDialog() {
  const { jwt } = useAuth();
  const { fetchInternalAgents, fetchExternalAgents } = useAgentDashboard();
  const { agent, isOpen, setIsOpen } = useAgentImportationThirdStepDialog();

  const { walletProvider } = useAppKitProvider("eip155");
  const eip1193Provider = walletProvider as Eip1193Provider;

  const [isSubmitting, setIsSubmitting] = React.useState(false);

  // Function to send TransactionSignature to the smart contract
  const sendUpdateSignerTransaction = async (
    nonce: string,
    tokenId: number,
    newSigner: string,
    signature: string,
    collectionAddress: string
  ) => {
    const tx = await updateAgentSigner(
      collectionAddress,
      tokenId,
      newSigner,
      nonce,
      signature,
      eip1193Provider
    );

    await tx.wait();
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Finish agent importation</DialogTitle>
          <DialogDescription>
            You don't known how to import an agent? See our documentation.
          </DialogDescription>
        </DialogHeader>
        <div className="flex flex-col space-y-8">
          <div className="flex flex-col space-y-4">
            <span>Your about to finish the importation of your agent.</span>
            <div className="flex flex-col space-y-1">
              <span>Make sure you :</span>
              <ul className="px-1 text-sm list-disc list-inside">
                <li>Created the database in your database server.</li>
                <li>
                  Created the role for the agent able to connect to the created
                  database
                </li>
                <li>Created the slack account for the agent</li>
              </ul>
            </div>
          </div>
          <div className="flex w-full justify-end">
            <Button
              type="button"
              loading={isSubmitting}
              onClick={async () => {
                setIsSubmitting(true);

                // Send TransactionSignature to the smart contract
                try {
                  const response = await importAgentThirdStep(
                    jwt!,
                    agent!.collectionAddress,
                    agent!.collectionTokenId!
                  );

                  if (response.status !== 200) {
                    toast("Failed to finish agent creation");
                    return;
                  }

                  // nonce is the agent id str without hyphens
                  const nonce = agent!.id.replace(/-/g, "");
                  const byte32Nonce = ethers.keccak256(strToBytes(nonce));

                  await sendUpdateSignerTransaction(
                    byte32Nonce,
                    agent!.collectionTokenId!,
                    response.data.signer,
                    response.data.signature,
                    agent!.collectionAddress
                  );

                  fetchInternalAgents();
                  fetchExternalAgents();
                  setIsOpen(false);

                  toast("Agent importation finished");
                } catch (error) {
                  console.error(error);
                  toast("Failed to send transaction to smart contract");
                } finally {
                  setIsSubmitting(false);
                }
              }}
            >
              {isSubmitting && <Loader2Icon className="animate-spin" />}
              Finish agent importation
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
