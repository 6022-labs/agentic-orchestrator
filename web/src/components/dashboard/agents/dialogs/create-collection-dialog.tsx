import React from "react"
import { toast } from "sonner"
import { Plus, TriangleAlert } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import CreateCollectionForm from "@/components/dashboard/agents/forms/create-collection-form"
import type { CreateCollectionFormSchemaType } from "@/components/dashboard/agents/forms/create-collection-form/create-collection-form-schema"
import { useAuth } from "@/contexts/auth-context"
import { useWallet } from "@/contexts/wallet-provider-context"
import type { FormCompletionResult } from "@/models/api/validation-error"
import {
  canCreateCollection,
  createCollection,
} from "@/queries/contracts/agent-collection-creator"
import { extractErrorMessage } from "@/lib/api-errors"

type Props = {
  onCreated: (created: {
    address: string | null
    name: string
  }) => Promise<void> | void
}

export default function CreateCollectionDialog({ onCreated }: Props) {
  const { user } = useAuth()
  const { wallet } = useWallet()

  const [open, setOpen] = React.useState(false)
  const [saving, setSaving] = React.useState(false)
  // null = probe unavailable; show no hint rather than guessing.
  const [allowed, setAllowed] = React.useState<boolean | null>(null)

  React.useEffect(() => {
    if (!open || !wallet || !user?.address) return
    let active = true
    canCreateCollection(user.address, wallet)
      .then((ok) => {
        if (active) setAllowed(ok)
      })
      .catch(() => {
        if (active) setAllowed(null)
      })
    return () => {
      active = false
    }
  }, [open, wallet, user?.address])

  const onFormComplete = async (
    values: CreateCollectionFormSchemaType
  ): Promise<FormCompletionResult<keyof CreateCollectionFormSchemaType>> => {
    if (!wallet) {
      toast.error("Connect a wallet first")
      return {}
    }

    setSaving(true)
    try {
      const address = await createCollection(
        values.name,
        values.symbol,
        values.description ?? "",
        values.admin,
        values.moderator,
        wallet
      )
      toast.success("Collection created")
      setOpen(false)
      await onCreated({ address, name: values.name })
      return {}
    } catch (error) {
      toast.error(extractErrorMessage(error, "Failed to create collection"))
      return {}
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button type="button" variant="outline" size="sm">
          <Plus />
          New collection
        </Button>
      </DialogTrigger>
      <DialogContent className="max-h-[85vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Create collection</DialogTitle>
          <DialogDescription>
            Deploys a new agent collection from your wallet and registers it in
            the on-chain manager.
          </DialogDescription>
        </DialogHeader>

        {allowed === false ? (
          <div className="flex items-start gap-2 rounded-md bg-amber-500/10 p-3 text-sm text-amber-600 dark:text-amber-500">
            <TriangleAlert className="mt-0.5 size-4 shrink-0" />
            <span>
              Your wallet can't create collections — it needs the minimum
              authorized-token balance or locked collateral.
            </span>
          </div>
        ) : null}

        <CreateCollectionForm onFormComplete={onFormComplete} />

        <DialogFooter>
          <Button
            type="submit"
            form="create-collection-form"
            loading={saving}
            disabled={allowed === false}
          >
            Create collection
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
