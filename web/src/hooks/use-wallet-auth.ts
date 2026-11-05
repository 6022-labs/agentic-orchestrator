import React from "react"
import { toast } from "sonner"
import { BrowserProvider } from "ethers"
import { useAppKit } from "@reown/appkit/react"

import { useAuth } from "@/contexts/auth-context"
import { useWallet } from "@/contexts/wallet-provider-context"
import { getAuthChallenge, verifyAuth } from "@/queries/api/auth.queries"
import { extractErrorMessage } from "@/lib/api-errors"

// SIWE (EIP-4361) for the gateway: connect a wallet (Reown AppKit), fetch the
// server-built sign-in message for the address, personal_sign it verbatim,
// then exchange it for a JWT.
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
      const network = await provider.getNetwork()

      const { data: challenge } = await getAuthChallenge(
        address,
        Number(network.chainId)
      )
      const signature = await signer.signMessage(challenge.message)

      const { data } = await verifyAuth(address, signature)
      signIn(data.token, data.refreshToken)
      toast.success("Signed in")
    } catch (error) {
      toast.error(extractErrorMessage(error, "Sign-in failed"))
    } finally {
      setIsAuthenticating(false)
    }
  }, [wallet, open, signIn])

  return { authenticate, isAuthenticating, isConnected }
}
