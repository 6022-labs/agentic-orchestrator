import React from "react";
import { toast } from "sonner";
import { Eip1193Provider } from "ethers";
import { MintProposal } from "@/models/mint-proposal";
import MintProposalCard from "./cards/MintProposalCard";
import { useAppKitProvider } from "@reown/appkit/react";
import { useUserModerator } from "@/contexts/user-moderator-context";
import {
  getMintProposals,
  mintFromProposal,
  refuseMintProposal,
} from "@/queries/contracts/agent-collection";

export default function MintProposals() {
  const [mintProposals, setMintProposals] = React.useState<
    Record<string, MintProposal[]>
  >({});

  const { walletProvider } = useAppKitProvider("eip155");
  const eip1193Provider = walletProvider as Eip1193Provider;

  const { collections } = useUserModerator();

  async function fetchMintProposals() {
    if (!collections || collections.length === 0) {
      setMintProposals({});
      return;
    }
    try {
      for (const collection of collections) {
        const proposals = await getMintProposals(
          collection.address,
          eip1193Provider
        );

        setMintProposals((prev) => ({
          ...prev,
          [collection.address]: proposals,
        }));
      }
    } catch (error) {
      console.error("Failed to fetch mint proposals:", error);
      toast.error("Failed to fetch mint proposals");
    }
  }

  React.useEffect(() => {
    fetchMintProposals();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [collections]);

  return (
    <div className="flex flex-col gap-6">
      {!Object.values(mintProposals).some((arr) => arr.length > 0) ? (
        <div className="text-center text-gray-500 py-8">
          No proposals to display
        </div>
      ) : (
        Object.entries(mintProposals)
          .filter(([, proposals]) => proposals.length > 0)
          .map(([collectionAddress, proposals]) => {
            const collection = collections.find(
              (c) => c.address === collectionAddress
            );

            return (
              <div key={collectionAddress}>
                <h2 className="text-xl font-semibold mb-4">
                  <div className="flex items-center gap-2">
                    <span>
                      Proposals for {collection?.name || collectionAddress}
                    </span>
                    <span className="text-gray-500 text-sm">
                      ({collectionAddress})
                    </span>
                  </div>
                </h2>
                <div className="grid grid-cols-1 gap-6">
                  {proposals.map((proposal) => (
                    <MintProposalCard
                      key={proposal.id}
                      proposal={proposal}
                      onAccept={async () => {
                        try {
                          await mintFromProposal(
                            collectionAddress,
                            proposal.id,
                            eip1193Provider
                          );

                          toast.success(
                            `Proposal ${proposal.id} accepted successfully`
                          );

                          await fetchMintProposals();
                        } catch {
                          toast.error(
                            `Failed to accept proposal ${proposal.id}`
                          );
                          return;
                        }
                      }}
                      onReject={async () => {
                        try {
                          await refuseMintProposal(
                            collectionAddress,
                            proposal.id,
                            eip1193Provider
                          );

                          toast.success(
                            `Proposal ${proposal.id} rejected successfully`
                          );

                          await fetchMintProposals();
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
