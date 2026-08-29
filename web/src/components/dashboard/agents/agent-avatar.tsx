import { Bot } from "lucide-react"

import { cn } from "@/lib/utils"
import IpfsImage from "@/components/ipfs-image"
import type { Agent } from "@/models/agent"

// Agent avatar resolved from the "default" (or first) IPFS image CID.
export default function AgentAvatar({
  agent,
  className,
}: {
  agent: Agent
  className?: string
}) {
  const cid = agent.images?.default ?? Object.values(agent.images ?? {})[0]

  return (
    <IpfsImage
      value={cid}
      alt={agent.name}
      className={cn("rounded-full object-cover", className)}
      fallback={
        <span
          className={cn(
            "flex items-center justify-center rounded-full bg-muted text-muted-foreground",
            className
          )}
        >
          <Bot className="size-1/2" />
        </span>
      }
    />
  )
}
