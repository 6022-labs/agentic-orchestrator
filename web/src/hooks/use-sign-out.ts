import * as React from "react"
import { useDisconnect } from "@reown/appkit/react"

import { useAuth } from "@/contexts/auth-context"

// useSignOut is the user's explicit sign-out; it also disconnects the wallet, as signing out of a dApp is expected to.
export const useSignOut = () => {
  const { signOutExplicitly } = useAuth()
  const { disconnect } = useDisconnect()

  return React.useCallback(() => {
    signOutExplicitly()
    disconnect().catch(() => undefined)
  }, [signOutExplicitly, disconnect])
}
