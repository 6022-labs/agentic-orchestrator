import React from "react";
import { toast } from "sonner";
import { Eip1193Provider } from "ethers";
import { useAppKitProvider } from "@reown/appkit/react";
import { useUserModerator } from "@/contexts/user-moderator-context";
import AddOrUpdateImageProposalCard from "./cards/AddOrUpdateImageProposalCard";
import { AddOrUpdateImageProposal } from "@/models/add-or-update-image-proposal";
import {
  addOrUpdateAgentImageFromProposal,
  getAddOrUpdateAgentImageProposals,
  refuseAddOrUpdateAgentImageProposal,
} from "@/queries/contracts/agent-collection";

export default function AddOrUpdateAgentImageProposals() {
  const [addOrUpdateAgentImageProposals, setAddOrUpdateAgentImageProposals] =
    React.useState<Record<string, AddOrUpdateImageProposal[]>>({});

  const { walletProvider } = useAppKitProvider("eip155");
  const eip1193Provider = walletProvider as Eip1193Provider;

  const { collections } = useUserModerator();

  async function fetchAddOrUpdateAgentImageProposals() {
    if (!collections || collections.length === 0) {
      setAddOrUpdateAgentImageProposals({});
      return;
    }
    try {
      for (const collection of collections) {
        const proposals = await getAddOrUpdateAgentImageProposals(
          collection.address,
          eip1193Provider
        );

        setAddOrUpdateAgentImageProposals((prev) => ({
          ...prev,
          [collection.address]: proposals,
        }));
      }
    } catch (error) {
      console.error(
        "Failed to fetch add or update agent image proposals:",
        error
      );
      toast.error("Failed to fetch add or update agent image proposals");
    }
  }

  React.useEffect(() => {
    fetchAddOrUpdateAgentImageProposals();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [collections]);

  return (
    <div className="flex flex-col gap-6">
      {!Object.values(addOrUpdateAgentImageProposals).some(
        (arr) => arr.length > 0
      ) ? (
        <div className="text-center text-gray-500 py-8">
          No proposals to display
        </div>
      ) : (
        Object.entries(addOrUpdateAgentImageProposals)
          .filter(([, proposals]) => proposals.length > 0)
          .map(([collectionAddress, proposals]) => {
            const collection = collections.find(
              (c) => c.address === collectionAddress
            );

            return (
              <div key={collectionAddress}>
                <h2 className="text-xl font-semibold mb-4">
                  Proposals for {collection?.name || collectionAddress} (
                  {collectionAddress})
                </h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                  {proposals.map((proposal) => (
                    <AddOrUpdateImageProposalCard
                      key={proposal.id}
                      proposal={proposal}
                      onAccept={async () => {
                        try {
                          await addOrUpdateAgentImageFromProposal(
                            collectionAddress,
                            proposal.id,
                            eip1193Provider
                          );

                          toast.success(
                            `Proposal ${proposal.id} accepted successfully`
                          );

                          await fetchAddOrUpdateAgentImageProposals();
                        } catch {
                          toast.error(
                            `Failed to accept proposal ${proposal.id}`
                          );
                          return;
                        }
                      }}
                      onReject={async () => {
                        try {
                          await refuseAddOrUpdateAgentImageProposal(
                            collectionAddress,
                            proposal.id,
                            eip1193Provider
                          );

                          toast.success(
                            `Proposal ${proposal.id} rejected successfully`
                          );

                          await fetchAddOrUpdateAgentImageProposals();
                        } catch {
                          toast.error(
                            `Failed to reject proposal ${proposal.id}`
                          );
                          return;
                        }
                      }}
                    />
                  ))}
                </div>
              </div>
            );
          })
      )}
    </div>
  );
}
