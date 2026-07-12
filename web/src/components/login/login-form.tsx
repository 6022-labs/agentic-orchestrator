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
  const { jwt } = useAuth()
  const { wallet, isConnected } = useWallet()
  const { authenticate, isAuthenticating } = useWalletAuth()
  const autoSigned = useRef(false)

  // Mirror the orchestrator's one-shot flow: once the wallet is connected and
  // its provider is ready, prompt the signature automatically instead of
  // requiring a second click.
  useEffect(() => {
    if (
      isConnected &&
      wallet &&
      !jwt &&
      !isAuthenticating &&
      !autoSigned.current
    ) {
      autoSigned.current = true
      authenticate()
    }
    if (!isConnected) {
      autoSigned.current = false
    }
  }, [isConnected, wallet, jwt, isAuthenticating, authenticate])

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
          onClick={authenticate}
        >
          Login
        </Button>
      </div>
    </form>
  )
}
