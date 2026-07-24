import React from "react"
import { toast } from "sonner"
import { ArrowRightLeft } from "lucide-react"

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
import TransferAgentForm from "@/components/dashboard/agents/forms/transfer-agent-form"
import type { TransferAgentFormSchemaType } from "@/components/dashboard/agents/forms/transfer-agent-form/transfer-agent-form-schema"
import { useWallet } from "@/contexts/wallet-provider-context"
import type { FormCompletionResult } from "@/models/api/validation-error"
import { transferFrom } from "@/queries/contracts/agent-collection"
import { extractErrorMessage } from "@/lib/api-errors"

type Props = {
  collectionAddress: string
  collectionTokenId: string
  agentName?: string
  onTransferred: () => Promise<void> | void
}

export default function TransferAgentDialog({
  collectionAddress,
  collectionTokenId,
  agentName,
  onTransferred,
}: Props) {
  const { wallet } = useWallet()

  const [open, setOpen] = React.useState(false)
  const [saving, setSaving] = React.useState(false)

  const onFormComplete = async (
    values: TransferAgentFormSchemaType
  ): Promise<FormCompletionResult<keyof TransferAgentFormSchemaType>> => {
    if (!wallet) {
      toast.error("Connect a wallet first")
      return {}
    }

    setSaving(true)
    try {
      await transferFrom(collectionAddress, collectionTokenId, values.to, wallet)
      toast.success("Agent transferred")
      setOpen(false)
      await onTransferred()
      return {}
    } catch (error) {
      toast.error(extractErrorMessage(error, "Failed to transfer the agent"))
      return {}
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="size-6"
          title="Transfer agent"
          aria-label="Transfer agent"
          onClick={(event) => event.stopPropagation()}
        >
          <ArrowRightLeft className="size-3.5" />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Transfer agent</DialogTitle>
          <DialogDescription>
            {agentName ? `Transfers "${agentName}"` : "Transfers this agent"} to
            another wallet. The transaction is signed by the current owner.
          </DialogDescription>
        </DialogHeader>

        <TransferAgentForm onFormComplete={onFormComplete} />

        <DialogFooter>
          <Button type="submit" form="transfer-agent-form" loading={saving}>
            Transfer
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
