import React from "react"
import { Loader2, RotateCcw, Unplug } from "lucide-react"

import { Button } from "@/components/ui/button"
import { loadAppConfig } from "@/configs/app-config"

type GateState = "loading" | "error" | "ready"

// Blocks rendering until the orchestrator gateway's /config is loaded — the app
// needs the chain id before AppKit initializes.
export default function ConfigGate({
  children,
}: {
  children: React.ReactNode
}) {
  const [state, setState] = React.useState<GateState>("loading")

  const load = React.useCallback(() => {
    setState("loading")
    loadAppConfig()
      .then(() => setState("ready"))
      .catch(() => setState("error"))
  }, [])

  React.useEffect(() => {
    load()
  }, [load])

  if (state === "ready") return children

  return (
    <div className="flex h-svh flex-col items-center justify-center gap-3">
      {state === "loading" ? (
        <>
          <Loader2 className="size-6 animate-spin text-muted-foreground" />
          <p className="text-sm text-muted-foreground">
            Connecting to the orchestrator…
          </p>
        </>
      ) : (
        <>
          <Unplug className="size-8 text-muted-foreground" />
          <p className="font-medium">Could not reach the orchestrator gateway</p>
          <p className="text-sm text-muted-foreground">
            Check that the orchestrator is running, then try again.
          </p>
          <Button variant="outline" onClick={load}>
            <RotateCcw />
            Retry
          </Button>
        </>
      )}
    </div>
  )
}
