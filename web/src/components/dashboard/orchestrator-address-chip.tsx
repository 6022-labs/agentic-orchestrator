import BlockiesRaw from "react-blockies"
import { Copy } from "lucide-react"
import { toast } from "sonner"

import { getAppConfig } from "@/configs/app-config"
import { truncateEthAddress } from "@/lib/utils"
import { Button } from "@/components/ui/button"

// react-blockies ships CJS; under this bundler the default import can arrive
// wrapped as { default }, so unwrap it to get the actual component.
const Blockies = ((BlockiesRaw as unknown as { default?: typeof BlockiesRaw })
  .default ?? BlockiesRaw) as typeof BlockiesRaw

export default function OrchestratorAddressChip() {
  const orchestratorAddress = getAppConfig().orchestratorAddress

  const copyOrchestratorAddress = async () => {
    await navigator.clipboard.writeText(orchestratorAddress)
    toast.success("Orchestrator address copied")
  }

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={copyOrchestratorAddress}
      title="Copy orchestrator address"
    >
      <span className="flex size-5 shrink-0 items-center justify-center overflow-hidden rounded-full">
        <Blockies seed={orchestratorAddress.toLowerCase()} size={8} scale={3} />
      </span>
      <span className="font-mono text-xs">
        {truncateEthAddress(orchestratorAddress)}
      </span>
      <Copy className="size-3.5" />
    </Button>
  )
}
