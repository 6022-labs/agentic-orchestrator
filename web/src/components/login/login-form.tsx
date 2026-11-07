import { useEffect, useRef } from "react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/contexts/auth-context"
import { useWallet } from "@/contexts/wallet-provider-context"
import { useWalletAuth } from "@/hooks/use-wallet-auth"

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"form">) {
  const { jwt, skipAutoSign, storedAddress, adoptStoredSession } = useAuth()
  const { wallet, address, isConnected } = useWallet()
  const { authenticate, isAuthenticating } = useWalletAuth()
  const autoSigned = useRef(false)
  const loginRequested = useRef(false)

  // Mirror the orchestrator's one-shot flow: once the wallet is connected and
  // its provider is ready, prompt the signature automatically instead of
  // requiring a second click.
  useEffect(() => {
    if (!isConnected) {
      autoSigned.current = false
      return
    }
    if (!wallet || jwt || isAuthenticating) return
    // Another tab already holds a session for this wallet: join it instead of signing again.
    if (
      address &&
      address.toLowerCase() === storedAddress?.toLowerCase() &&
      adoptStoredSession(address)
    )
      return
    if (autoSigned.current) return
    if (skipAutoSign && !loginRequested.current) return
    autoSigned.current = true
    authenticate()
  }, [
    isConnected,
    wallet,
    address,
    jwt,
    isAuthenticating,
    authenticate,
    skipAutoSign,
    storedAddress,
    adoptStoredSession,
  ])

  const login = () => {
    loginRequested.current = true
    if (wallet) autoSigned.current = true
    authenticate()
  }

  return (
    <form
      className={cn("flex flex-col gap-6", className)}
      onSubmit={(e) => e.preventDefault()}
      {...props}
    >
      <div className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-2xl font-bold">Login to your account</h1>
        <p className="text-sm text-balance text-muted-foreground">
          Use your web3 wallet to login
        </p>
      </div>
      <div className="grid gap-6">
        <Button
          type="button"
          className="w-full"
          loading={isAuthenticating}
          onClick={login}
        >
          Login
        </Button>
      </div>
    </form>
  )
}
