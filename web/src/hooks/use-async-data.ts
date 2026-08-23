import React from "react"
import { toast } from "sonner"

import { extractErrorMessage } from "@/lib/api-errors"

// One independent data source: its own loading flag, its own error toast — so a
// failing source never blocks or blanks the widgets fed by the other ones.
// `load` must be memoized by the caller (useCallback).
export function useAsyncData<T>(load: () => Promise<T>, errorMessage: string) {
  const [data, setData] = React.useState<T | undefined>(undefined)
  const [loading, setLoading] = React.useState(true)
  const [nonce, setNonce] = React.useState(0)

  const reload = React.useCallback(() => setNonce((current) => current + 1), [])

  React.useEffect(() => {
    let cancelled = false

    setLoading(true)
    load()
      .then((value) => {
        if (!cancelled) setData(value)
      })
      .catch((error) => {
        if (cancelled) return
        toast.error(extractErrorMessage(error, errorMessage))
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [load, errorMessage, nonce])

  return { data, loading, reload }
}
