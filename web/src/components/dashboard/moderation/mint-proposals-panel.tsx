import { Inbox } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import ProposalAddress from "@/components/dashboard/moderation/proposal-address"
import ProposalsPanel from "@/components/dashboard/moderation/proposals-panel"
import { ipfsImageUrl } from "@/lib/ipfs"
import type { AgentCollection } from "@/models/agent-collection"
import type { MintProposal } from "@/models/mint-proposal"
import {
  getMintProposals,
  mintFromProposal,
  refuseMintProposal,
} from "@/queries/contracts/agent-collection"

export default function MintProposalsPanel({
  collections,
}: {
  collections: AgentCollection[]
}) {
  return (
    <ProposalsPanel<MintProposal>
      title="Mint proposals"
      description="Pending agent mints awaiting moderator approval."
      emptyIcon={Inbox}
      emptyTitle="No mint proposals"
      collections={collections}
      loadErrorMessage="Failed to load mint proposals"
      list={getMintProposals}
      accept={mintFromProposal}
      reject={refuseMintProposal}
      getImage={(proposal) => ({
        url: ipfsImageUrl(proposal.images[0]?.value) ?? undefined,
        alt: proposal.name,
      })}
      renderBody={(proposal) => (
        <>
          <div className="flex items-center gap-2">
            <span className="truncate font-medium">{proposal.name}</span>
            <Badge variant="secondary">{proposal.role}</Badge>
          </div>
          {proposal.mission ? (
            <p className="text-sm text-muted-foreground">{proposal.mission}</p>
          ) : null}
          <ProposalAddress label="Owner" address={proposal.owner} />
        </>
      )}
      rejectDialog={(proposal) => ({
        title: "Reject this mint proposal?",
        description: `The proposal for "${proposal.name}" is refused on-chain. This cannot be undone.`,
      })}
    />
  )
}
