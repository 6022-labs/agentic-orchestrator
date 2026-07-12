import { ExternalLink } from "lucide-react"

import { Button } from "@/components/ui/button"
import { config } from "@/configs/env"
import { truncateEthAddress } from "@/lib/utils"

// Address row with a block-explorer link, shared by the proposal cards.
export default function ProposalAddress({
  label,
  address,
}: {
  label: string
  address: string
}) {
  return (
    <div className="flex items-center gap-1 text-sm">
      <span className="text-muted-foreground">{label}</span>
      <span className="font-mono">{truncateEthAddress(address, 8, 6)}</span>
      <Button variant="ghost" size="icon" className="size-6" asChild>
        <a
          href={`${config.scannerBaseUrl}/address/${address}`}
          target="_blank"
          rel="noreferrer"
        >
          <ExternalLink className="size-3.5" />
        </a>
      </Button>
    </div>
  )
}
