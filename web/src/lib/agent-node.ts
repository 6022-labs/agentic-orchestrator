import { toast } from "sonner"

import { extractErrorMessage } from "@/lib/api-errors"
import { getAgentUrl } from "@/queries/api/agents.queries"

// Opens a tab synchronously so popup blockers allow it, then points it at the
// node URL resolved from the agent's ENS url record.
export const openAgentNode = (
  collectionAddress: string,
  collectionTokenId: string
) => {
  const win = window.open("", "_blank")
  getAgentUrl(collectionAddress, collectionTokenId)
    .then(({ data }) => {
      if (win) win.location.href = data.url
    })
    .catch((error) => {
      win?.close()
      toast.error(extractErrorMessage(error, "Failed to resolve agent URL"))
    })
}
