import React from "react"
import { toast } from "sonner"
import { BrowserProvider, type Eip1193Provider } from "ethers"

import { useWallet } from "@/contexts/wallet-provider-context"
import { verifyWalletForBridge } from "@/queries/api/wallet.queries"
import { extractErrorMessage } from "@/lib/api-errors"

// Bridge wallet pairing: sign the raw challenge from the platform bridge link
// (e.g. Slack) and post it back so the gateway binds the wallet to the user.
export const useBridgeVerification = (
  platform: string | null,
  challenge: string | null
) => {
  const { wallet } = useWallet()

  const [isProcessing, setIsProcessing] = React.useState(false)
  const [errorMessage, setErrorMessage] = React.useState<string | null>(null)
  const [isSuccess, setIsSuccess] = React.useState(false)

  React.useEffect(() => {
    if (!challenge || !platform) {
      setErrorMessage("Missing required challenge information.")
    } else {
      setErrorMessage(null)
    }
  }, [challenge, platform])

  const signChallenge = React.useCallback(async () => {
    if (!challenge || !platform) {
      setErrorMessage("Missing required challenge information.")
      return
    }

    setIsProcessing(true)
    setErrorMessage(null)

    try {
      const provider = new BrowserProvider(wallet as Eip1193Provider)
      const signer = await provider.getSigner()
      const signature = await signer.signMessage(challenge)

      await verifyWalletForBridge(platform, challenge, signature)

      toast.success("Challenge signed successfully.")
      setIsSuccess(true)
    } catch (error) {
      const message = extractErrorMessage(
        error,
        "An error occurred while signing the challenge."
      )
      toast.error(message)
      setErrorMessage(message)
    } finally {
      setIsProcessing(false)
    }
  }, [challenge, platform, wallet])

  return { signChallenge, isProcessing, errorMessage, isSuccess }
}
