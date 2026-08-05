import React from "react"
import { toast } from "sonner"
import type { LucideIcon } from "lucide-react"
import type { Eip1193Provider } from "ethers"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import ConfirmDialog from "@/components/dashboard/confirm-dialog"
import EmptyState from "@/components/dashboard/empty-state"
import { extractErrorMessage } from "@/lib/api-errors"
import type { AgentCollection } from "@/models/agent-collection"
import { useWallet } from "@/contexts/wallet-provider-context"

type ProposalsPanelProps<TProposal extends { id: number }> = {
  title: string
  description: string
  emptyIcon: LucideIcon
  emptyTitle: string
  collections: AgentCollection[]
  loadErrorMessage: string
  list: (
    collectionAddress: string,
    wallet: Eip1193Provider
  ) => Promise<TProposal[]>
  accept: (
    collectionAddress: string,
    proposalId: number,
    wallet: Eip1193Provider
  ) => Promise<unknown>
  reject: (
    collectionAddress: string,
    proposalId: number,
    wallet: Eip1193Provider
  ) => Promise<unknown>
  getImage: (proposal: TProposal) => { url?: string; alt: string }
  renderBody: (proposal: TProposal) => React.ReactNode
  rejectDialog: (proposal: TProposal) => { title: string; description: string }
}

// ProposalsPanel is the shared moderation card: it loads on-chain proposals per
// collection and drives the accept/reject flow; each proposal type only
// provides its queries and row content.
export default function ProposalsPanel<TProposal extends { id: number }>({
  title,
  description,
  emptyIcon,
  emptyTitle,
  collections,
  loadErrorMessage,
  list,
  accept,
  reject,
  getImage,
  renderBody,
  rejectDialog,
}: ProposalsPanelProps<TProposal>) {
  const { wallet } = useWallet()
  const [loading, setLoading] = React.useState(true)
  const [proposalsByCollection, setProposalsByCollection] = React.useState<
    Record<string, TProposal[]>
  >({})
  const [actingOn, setActingOn] = React.useState<string | null>(null)

  const load = React.useCallback(async () => {
    if (!wallet || collections.length === 0) {
      setProposalsByCollection({})
      setLoading(false)
      return
    }

    setLoading(true)
    try {
      const entries = await Promise.all(
        collections.map(
          async (collection) =>
            [collection.address, await list(collection.address, wallet)] as const
        )
      )
      setProposalsByCollection(Object.fromEntries(entries))
    } catch (error) {
      toast.error(extractErrorMessage(error, loadErrorMessage))
    } finally {
      setLoading(false)
    }
  }, [wallet, collections, list, loadErrorMessage])

  React.useEffect(() => {
    void load()
  }, [load])

  const act = async (
    collectionAddress: string,
    proposal: TProposal,
    isAccepted: boolean
  ) => {
    if (!wallet) return

    setActingOn(
      `${collectionAddress}/${proposal.id}:${isAccepted ? "accept" : "reject"}`
    )
    try {
      if (isAccepted) {
        await accept(collectionAddress, proposal.id, wallet)
      } else {
        await reject(collectionAddress, proposal.id, wallet)
      }
      toast.success(
        `Proposal ${isAccepted ? "accepted" : "rejected"} successfully`
      )
      await load()
    } catch (error) {
      toast.error(
        extractErrorMessage(
          error,
          `Failed to ${isAccepted ? "accept" : "reject"} the proposal`
        )
      )
    } finally {
      setActingOn(null)
    }
  }

  const hasProposals = Object.values(proposalsByCollection).some(
    (proposals) => proposals.length > 0
  )

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-6">
        {loading ? (
          <div className="space-y-2">
            <Skeleton className="h-24 w-full" />
            <Skeleton className="h-24 w-full" />
          </div>
        ) : !hasProposals ? (
          <EmptyState
            icon={emptyIcon}
            title={emptyTitle}
            description="Nothing is waiting for moderation right now."
          />
        ) : (
          Object.entries(proposalsByCollection)
            .filter(([, proposals]) => proposals.length > 0)
            .map(([collectionAddress, proposals]) => {
              const collection = collections.find(
                (c) =>
                  c.address.toLowerCase() === collectionAddress.toLowerCase()
              )

              return (
                <div key={collectionAddress} className="flex flex-col gap-3">
                  <h3 className="text-sm font-semibold">
                    {collection?.name || collectionAddress}
                  </h3>
                  {proposals.map((proposal) => {
                    const actionKey = `${collectionAddress}/${proposal.id}`
                    const image = getImage(proposal)
                    const dialog = rejectDialog(proposal)

                    return (
                      <div
                        key={actionKey}
                        className="flex flex-col gap-4 rounded-xl border p-4 sm:flex-row"
                      >
                        {image.url ? (
                          <img
                            src={image.url}
                            alt={image.alt}
                            className="size-24 shrink-0 rounded-lg object-cover"
                          />
                        ) : null}
                        <div className="flex min-w-0 flex-1 flex-col gap-2">
                          {renderBody(proposal)}
                        </div>
                        <div className="flex shrink-0 items-center justify-end gap-2 sm:flex-col sm:justify-center">
                          <Button
                            size="sm"
                            className="w-28"
                            loading={actingOn === `${actionKey}:accept`}
                            disabled={actingOn !== null}
                            onClick={() =>
                              void act(collectionAddress, proposal, true)
                            }
                          >
                            Accept
                          </Button>
                          <ConfirmDialog
                            title={dialog.title}
                            description={dialog.description}
                            confirmLabel="Reject"
                            onConfirm={() =>
                              void act(collectionAddress, proposal, false)
                            }
                          >
                            <Button
                              variant="outline"
                              size="sm"
                              className="w-28"
                              loading={actingOn === `${actionKey}:reject`}
                              disabled={actingOn !== null}
                            >
                              Reject
                            </Button>
                          </ConfirmDialog>
                        </div>
                      </div>
                    )
                  })}
                </div>
              )
            })
        )}
      </CardContent>
    </Card>
  )
}
