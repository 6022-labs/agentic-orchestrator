import { ImageIcon } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import ProposalsPanel from "@/components/dashboard/moderation/proposals-panel"
import type { AddOrUpdateImageProposal } from "@/models/add-or-update-image-proposal"
import type { AgentCollection } from "@/models/agent-collection"
import {
  addOrUpdateAgentImageFromProposal,
  getAddOrUpdateAgentImageProposals,
  refuseAddOrUpdateAgentImageProposal,
} from "@/queries/contracts/agent-collection"

export default function ImageProposalsPanel({
  collections,
}: {
  collections: AgentCollection[]
}) {
  return (
    <ProposalsPanel<AddOrUpdateImageProposal>
      title="Agent image proposals"
      description="Pending image additions or updates awaiting moderator approval."
      emptyIcon={ImageIcon}
      emptyTitle="No image proposals"
      collections={collections}
      loadErrorMessage="Failed to load image proposals"
      list={getAddOrUpdateAgentImageProposals}
      accept={addOrUpdateAgentImageFromProposal}
      reject={refuseAddOrUpdateAgentImageProposal}
      getImage={(proposal) => ({
        value: proposal.image?.value,
        alt: proposal.image.key,
      })}
      renderBody={(proposal) => (
        <div className="flex items-center gap-2">
          <span className="truncate font-medium">{proposal.image.key}</span>
          <Badge variant="outline">Token #{proposal.tokenId?.toString()}</Badge>
        </div>
      )}
      rejectDialog={(proposal) => ({
        title: "Reject this image proposal?",
        description: `The "${proposal.image.key}" image proposal is refused on-chain. This cannot be undone.`,
      })}
    />
  )
}
