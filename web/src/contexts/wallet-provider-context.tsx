/* eslint-disable react-refresh/only-export-components */
import * as React from "react"
import { toast } from "sonner"
import type { Eip1193Provider } from "ethers"
import { useAppKitAccount, useAppKitProvider } from "@reown/appkit/react"

import { useAuth } from "./auth-context"

type WalletProviderContextType = {
  wallet: Eip1193Provider | null
  address?: string
  isConnected: boolean
}

const WalletProviderContext = React.createContext<WalletProviderContextType>(
  {} as WalletProviderContextType
)

export const WalletProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { walletProvider } = useAppKitProvider("eip155")
  const eip1193Provider = (walletProvider as Eip1193Provider) ?? null
  const { address, isConnected } = useAppKitAccount()
  const { user, releaseSession } = useAuth()
  const wasConnected = React.useRef(false)

  // A session is bound to the address that signed in. If the wallet switches
  // accounts, drop the JWT so the user re-authenticates as the new address.
  React.useEffect(() => {
    if (
      user &&
      address &&
      address.toLowerCase() !== user.address.toLowerCase()
    ) {
      if (releaseSession())
        toast.info("Wallet account changed. Please sign in again.")
    }
  }, [address, user, releaseSession])

  // Disconnecting the wallet ends the session (mirrors the orchestrator UX).
  // Guarded on a prior connection so a page refresh doesn't drop the JWT before
  // the wallet has had a chance to reconnect.
  React.useEffect(() => {
    if (wasConnected.current && !isConnected && user) {
      releaseSession()
    }
    wasConnected.current = isConnected
  }, [isConnected, user, releaseSession])

  return (
    <WalletProviderContext.Provider
      value={{ wallet: eip1193Provider, address, isConnected }}
    >
      {children}
    </WalletProviderContext.Provider>
  )
}

export const useWallet = () => {
  const context = React.useContext(WalletProviderContext)
  if (!context)
    throw new Error("useWallet must be used within a WalletProvider")
  return context
}
