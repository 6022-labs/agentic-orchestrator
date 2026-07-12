import { Link } from "react-router-dom"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { useBridgeVerification } from "@/components/verify-wallet/use-bridge-verification"

export default function VerifyWalletSection({
  platform,
  challenge,
}: {
  platform: string | null
  challenge: string | null
}) {
  const { signChallenge, isProcessing, errorMessage, isSuccess } =
    useBridgeVerification(platform, challenge)

  if (isSuccess) {
    return (
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Verification complete</CardTitle>
          <CardDescription>
            Your wallet challenge has been signed successfully. You can now
            close this page or continue to the dashboard.
          </CardDescription>
        </CardHeader>
        <CardContent className="flex">
          <Button asChild>
            <Link to="/dashboard/overview">Go to Dashboard</Link>
          </Button>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="w-full max-w-lg">
      <CardHeader>
        <CardTitle>Wallet challenge verification</CardTitle>
        <CardDescription>
          {platform ? (
            <>
              <span>Sign the challenge to continue with </span>
              <span className="font-bold">{platform.toLowerCase()}</span>
            </>
          ) : (
            <span>Sign the challenge with your wallet to continue.</span>
          )}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div>
          <p className="text-sm font-medium text-muted-foreground">Challenge</p>
          <pre className="mt-2 max-h-48 overflow-auto rounded-md bg-muted p-3 text-sm">
            {challenge ?? "No challenge provided."}
          </pre>
        </div>

        {errorMessage ? (
          <div className="rounded-md border border-destructive/40 bg-destructive/10 p-3 text-sm text-destructive">
            {errorMessage}
          </div>
        ) : null}

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
          <Button
            type="button"
            className="w-full sm:w-auto"
            onClick={() => void signChallenge()}
            loading={isProcessing}
            disabled={isProcessing || !challenge || !platform}
          >
            Sign challenge
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
