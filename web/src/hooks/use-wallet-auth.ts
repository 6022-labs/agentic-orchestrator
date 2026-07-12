import React from "react"
import { toast } from "sonner"
import { BrowserProvider } from "ethers"
import { useAppKit } from "@reown/appkit/react"

import { useAuth } from "@/contexts/auth-context"
import { useWallet } from "@/contexts/wallet-provider-context"
import { getAuthChallenge, verifyAuth } from "@/queries/api/auth.queries"
import { extractErrorMessage } from "@/lib/api-errors"

// SIWE for the agent-node gateway: connect a wallet (Reown AppKit), then
// challenge -> personal_sign(nonce) -> verify. Unlike a full SIWE message, the
// gateway recovers the signer from the raw nonce, so we sign the nonce directly.
export const useWalletAuth = () => {
  const { open } = useAppKit()
  const { wallet, isConnected } = useWallet()
  const { signIn } = useAuth()
  const [isAuthenticating, setIsAuthenticating] = React.useState(false)

  const authenticate = React.useCallback(async () => {
    if (!wallet) {
      // No connected wallet yet — open the AppKit modal to connect first.
      await open()
      return
    }

    setIsAuthenticating(true)
    try {
      const provider = new BrowserProvider(wallet)
      const signer = await provider.getSigner()
      const address = await signer.getAddress()

      const { data: challenge } = await getAuthChallenge()
      const signature = await signer.signMessage(challenge.nonce)

      const { data } = await verifyAuth(address, signature, challenge.nonce)
      signIn(data.token)
      toast.success("Signed in")
    } catch (error) {
      toast.error(extractErrorMessage(error, "Sign-in failed"))
    } finally {
      setIsAuthenticating(false)
    }
  }, [wallet, open, signIn])

  return { authenticate, isAuthenticating, isConnected }
}
