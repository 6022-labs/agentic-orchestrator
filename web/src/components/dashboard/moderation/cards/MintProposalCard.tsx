import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { MintProposal } from "@/models/mint-proposal";
import { CopyButton } from "@/components/ui/copy-button";
import { Card, CardContent } from "@/components/ui/card";
import { ViewOnScannerLink } from "@/components/ui/view-on-scanner-link";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";

export type MintProposalCardProps = {
  proposal: MintProposal;
  className?: string;
  onAccept: () => void;
  onReject: () => void;
};

export default function MintProposalCard({
  proposal,
  className,
  onAccept,
  onReject,
}: MintProposalCardProps) {
  return (
    <Card className={cn("flex flex-row py-0", className)}>
      <div className="w-56 flex-shrink-0 relative">
        <Carousel className="w-full">
          <CarouselContent>
            {proposal.images.map((img) => (
              <CarouselItem key={img.key} className="relative h-56 w-full">
                <img
                  src={import.meta.env.VITE_IPFS_GATEWAY + img.value}
                  alt={img.key}
                  className="h-full w-full rounded-l-xl object-contain"
                />
              </CarouselItem>
            ))}
          </CarouselContent>
          {proposal.images.length > 1 && (
            <>
              <CarouselPrevious className="absolute top-auto bottom-1 left-1/3 z-10" />
              <CarouselNext className="absolute top-auto bottom-1 right-1/3 z-10" />
            </>
          )}
        </Carousel>
      </div>
      <div className="flex flex-col py-4 w-full justify-center">
        <CardContent className="flex w-full justify-between h-full items-center">
          <div className="flex flex-col gap-2 justify-between h-full">
            <div className="flex flex-col gap-2">
              <div className="flex flex-row gap-1 items-center">
                <span className="text-lg font-semibold">{proposal.name}</span>
                <span>({proposal.role})</span>
              </div>
              <div className="flex flex-row gap-1">
                <span className="font-bold">Mission:</span>
                <span>{proposal.mission}</span>
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
              <div className="flex flex-row items-center gap-1">
                <span className="font-bold">Owner:</span>
                <span>{proposal.owner}</span>
                <CopyButton className="w-7 h-7" text={proposal.owner} />
                <ViewOnScannerLink
                  className="w-7 h-7"
                  href={
                    import.meta.env.VITE_SCANNER_BASE_URL +
                    "/address/" +
                    proposal.owner
                  }
                />
              </div>
              <div className="flex flex-row items-center gap-1">
                <span className="font-bold">Signer:</span>
                <span>{proposal.signer}</span>
                <CopyButton className="w-7 h-7" text={proposal.signer} />
                <ViewOnScannerLink
                  className="w-7 h-7"
                  href={
                    import.meta.env.VITE_SCANNER_BASE_URL +
                    "/address/" +
                    proposal.signer
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
