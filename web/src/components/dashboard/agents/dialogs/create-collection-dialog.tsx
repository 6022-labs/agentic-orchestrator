import React from "react"
import { toast } from "sonner"
import { Plus } from "lucide-react"

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
import { useWallet } from "@/contexts/wallet-provider-context"
import type { FormCompletionResult } from "@/models/api/validation-error"
import { createCollection } from "@/queries/contracts/agent-collection-creator"
import { extractErrorMessage } from "@/lib/api-errors"

type Props = {
  onCreated: (created: {
    address: string | null
    name: string
  }) => Promise<void> | void
}

export default function CreateCollectionDialog({ onCreated }: Props) {
  const { wallet } = useWallet()

  const [open, setOpen] = React.useState(false)
  const [saving, setSaving] = React.useState(false)

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

        <CreateCollectionForm onFormComplete={onFormComplete} />

        <DialogFooter>
          <Button
            type="submit"
            form="create-collection-form"
            loading={saving}
          >
            Create collection
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
