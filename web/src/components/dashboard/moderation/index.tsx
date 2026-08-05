import { Gavel, Wallet } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import EmptyState from "@/components/dashboard/empty-state"
import ImageProposalsPanel from "@/components/dashboard/moderation/image-proposals-panel"
import MintProposalsPanel from "@/components/dashboard/moderation/mint-proposals-panel"
import { useModeratorCollections } from "@/components/dashboard/moderation/use-moderator-collections"
import { useWallet } from "@/contexts/wallet-provider-context"
import { useAppKit } from "@reown/appkit/react"

export default function ModerationSection() {
  const { open } = useAppKit()
  const { wallet, isConnected } = useWallet()
  const { collections, loading } = useModeratorCollections()

  if (!isConnected || !wallet) {
    return (
      <EmptyState
        icon={Wallet}
        title="Connect a wallet"
        description="Moderation runs on-chain, so a connected wallet is required."
        action={
          <Button variant="outline" onClick={() => open()}>
            Connect wallet
          </Button>
        }
      />
    )
  }

  if (loading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-40 w-full" />
        <Skeleton className="h-40 w-full" />
      </div>
    )
  }

  if (collections.length === 0) {
    return (
      <EmptyState
        icon={Gavel}
        title="No moderated collections"
        description="The connected wallet holds the moderator role on no collection."
      />
    )
  }

  return (
    <div className="flex flex-col gap-6">
      <MintProposalsPanel collections={collections} />
      <ImageProposalsPanel collections={collections} />
    </div>
  )
}
