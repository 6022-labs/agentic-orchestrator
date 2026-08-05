import { useSearchParams } from "react-router-dom"

import VerifyWalletSection from "@/components/verify-wallet"

export default function VerifyWalletPage() {
  const [searchParams] = useSearchParams()

  return (
    <div className="flex min-h-svh items-center justify-center bg-background p-4">
      <VerifyWalletSection
        platform={searchParams.get("platform")}
        challenge={searchParams.get("challenge")}
      />
    </div>
  )
}
