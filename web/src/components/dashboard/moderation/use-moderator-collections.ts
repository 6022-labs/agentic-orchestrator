import React from "react"
import { toast } from "sonner"

import { extractErrorMessage } from "@/lib/api-errors"
import type { AgentCollection } from "@/models/agent-collection"
import { listAgentCollections } from "@/queries/api/collections.queries"
import { moderatorOf } from "@/queries/contracts/agent-collections-manager"
import { useWallet } from "@/contexts/wallet-provider-context"

// Collections the connected wallet moderates: the on-chain manager's
// moderatorOf() filtered against the gateway's collection list.
export function useModeratorCollections() {
  const { wallet, address } = useWallet()
  const [collections, setCollections] = React.useState<AgentCollection[]>([])
  const [loading, setLoading] = React.useState(true)

  const load = React.useCallback(async () => {
    if (!wallet || !address) {
      setCollections([])
      setLoading(false)
      return
    }

    setLoading(true)
    try {
      const [{ data: allCollections }, moderatedAddresses] = await Promise.all([
        listAgentCollections(),
        moderatorOf(wallet),
      ])

      const moderated = new Set(
        (moderatedAddresses ?? []).map((a) => a.toLowerCase())
      )
      setCollections(
        (allCollections ?? []).filter((collection) =>
          moderated.has(collection.address.toLowerCase())
        )
      )
    } catch (error) {
      toast.error(
        extractErrorMessage(error, "Failed to load moderated collections")
      )
      setCollections([])
    } finally {
      setLoading(false)
    }
  }, [wallet, address])

  React.useEffect(() => {
    void load()
  }, [load])

  return { collections, loading, reload: load }
}
