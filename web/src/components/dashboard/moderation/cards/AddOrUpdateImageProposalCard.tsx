import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { CopyButton } from "@/components/ui/copy-button";
import { ViewOnScannerLink } from "@/components/ui/view-on-scanner-link";
import { AddOrUpdateImageProposal } from "@/models/add-or-update-image-proposal";

export type AddOrUpdateImageProposalCardProps = {
  proposal: AddOrUpdateImageProposal;
  className?: string;
  onAccept: () => Promise<void>;
  onReject: () => Promise<void>;
};

export default function AddOrUpdateImageProposalCard({
  proposal,
  className,
  onAccept,
  onReject,
}: AddOrUpdateImageProposalCardProps) {
  return (
    <Card className={cn("flex flex-row py-0", className)}>
      <div className="w-56 flex-shrink-0 relative">
        {proposal.image?.value && (
          <img
            src={import.meta.env.VITE_IPFS_GATEWAY + proposal.image.value}
            alt={proposal.image.key}
            className="h-56 w-full rounded-l-xl"
          />
        )}
      </div>
      <div className="flex flex-col py-4 w-full justify-center">
        <CardContent className="flex w-full justify-between h-full items-center">
          <div className="flex flex-col gap-2 justify-between h-full">
            <div className="flex flex-col gap-2">
              <div className="flex flex-row gap-1 items-center">
                <span className="text-lg font-semibold">
                  {proposal.image.key}
                </span>
                <span>(Image Key)</span>
              </div>
              <div className="flex flex-row gap-1">
                <span className="font-bold">Token ID:</span>
                <span>{proposal.tokenId?.toString()}</span>
              </div>
            </div>
            <div className="flex flex-col">
              <div className="flex flex-row items-center gap-1">
                <span className="font-bold">Collection Address:</span>
                <span>{proposal.collectionAddress}</span>
                <CopyButton
                  className="w-7 h-7"
                  text={proposal.collectionAddress}
                />
                <ViewOnScannerLink
                  className="w-7 h-7"
                  href={
                    import.meta.env.VITE_SCANNER_BASE_URL +
                    "/address/" +
                    proposal.collectionAddress
                  }
                />
              </div>
            </div>
          </div>
          <div className="flex flex-col gap-4 items-center justify-center h-full">
            <Button onClick={onAccept}>Accept</Button>
            <Button variant="outline" onClick={onReject}>
              Reject
            </Button>
          </div>
        </CardContent>
      </div>
    </Card>
  );
}
