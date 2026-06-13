import React from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Link, useSearchParams } from "react-router-dom";
import { BrowserProvider, Eip1193Provider } from "ethers";
import { useWallet } from "@/contexts/wallet-provider-context";
import { verifyWalletForBridge } from "@/queries/api/wallet.queries";
import {
  Card,
  CardTitle,
  CardHeader,
  CardContent,
  CardDescription,
} from "@/components/ui/card";

export default function VerifyWallet() {
  const [searchParams] = useSearchParams();
  const platform = searchParams.get("platform");
  const challenge = searchParams.get("challenge");

  const { wallet } = useWallet();

  const [isProcessing, setIsProcessing] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState<string | null>(null);

  const [isSuccess, setIsSuccess] = React.useState(false);

  React.useEffect(() => {
    if (!challenge || !platform) {
      setErrorMessage("Missing required challenge information.");
    } else {
      setErrorMessage(null);
    }
  }, [challenge, platform]);

  const handleSignChallenge = React.useCallback(async () => {
    if (!challenge || !platform) {
      setErrorMessage("Missing required challenge information.");
      return;
    }

    setIsProcessing(true);
    setErrorMessage(null);

    try {
      const provider = new BrowserProvider(wallet as Eip1193Provider);
      const signer = await provider.getSigner();
      const signature = await signer.signMessage(challenge);

      const response = await verifyWalletForBridge(
        platform,
        challenge,
        signature
      );

      if (response.status !== 200) {
        throw new Error("Failed to verify wallet challenge.");
      }

      toast.success("Challenge signed successfully.");
      setIsSuccess(true);
    } catch (error) {
      console.error(error);

      toast.error("An error occurred while signing the challenge.");
      setErrorMessage("An error occurred while signing the challenge.");
    } finally {
      setIsProcessing(false);
    }
  }, [challenge, platform, wallet]);

  return (
    <div className="min-h-svh bg-background flex items-center justify-center p-4">
      <Card className="max-w-lg w-full">
        {isSuccess ? (
          <>
            <CardHeader>
              <CardTitle>Verification complete</CardTitle>
              <CardDescription>
                Your wallet challenge has been signed successfully. You can now
                close this page or continue to the dashboard.
              </CardDescription>
            </CardHeader>
            <CardContent className="flex">
              <Button variant="default">
                <Link to="/dashboard/overview">Go to Dashboard</Link>
              </Button>
            </CardContent>
          </>
        ) : (
          <>
            <CardHeader>
              <CardTitle>Wallet Challenge Verification</CardTitle>
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
                <p className="text-sm font-medium text-muted-foreground">
                  Challenge
                </p>
                <pre className="mt-2 max-h-48 overflow-auto rounded-md bg-muted p-3 text-sm">
                  {challenge ?? "No challenge provided."}
                </pre>
              </div>

              {errorMessage && (
                <div className="rounded-md border border-destructive/40 bg-destructive/10 p-3 text-sm text-destructive">
                  {errorMessage}
                </div>
              )}

              <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
                <Button
                  type="button"
                  className="w-full sm:w-auto"
                  onClick={handleSignChallenge}
                  loading={isProcessing}
                  disabled={
                    isProcessing || !challenge || !platform || isSuccess
                  }
                >
                  {isProcessing ? "Processing..." : "Sign challenge"}
                </Button>
              </div>
            </CardContent>
          </>
        )}
      </Card>
    </div>
  );
}
