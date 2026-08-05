import React from "react"
import { Unplug } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import EmptyState from "@/components/dashboard/empty-state"
import type { AgentCard } from "@/models/agent-card"
import { getAgentCard } from "@/queries/api/agents.queries"

function CardField({
  label,
  value,
}: {
  label: string
  value?: React.ReactNode
}) {
  return (
    <div className="flex flex-col gap-1">
      <span className="text-xs text-muted-foreground">{label}</span>
      <span className="text-sm break-all">{value || "—"}</span>
    </div>
  )
}

// Live 6022 well-known card fetched through the gateway; the agent's node may
// be unreachable, in which case a graceful fallback is shown.
export default function AgentCardPanel({
  collectionAddress,
  collectionTokenId,
}: {
  collectionAddress: string
  collectionTokenId: string
}) {
  const [loading, setLoading] = React.useState(true)
  const [card, setCard] = React.useState<AgentCard | null>(null)
  const [unreachable, setUnreachable] = React.useState(false)

  const load = React.useCallback(async () => {
    setLoading(true)
    setUnreachable(false)
    try {
      const { data } = await getAgentCard(collectionAddress, collectionTokenId)
      setCard(data)
    } catch {
      // 404/502 both mean the agent's node did not answer with a card.
      setCard(null)
      setUnreachable(true)
    } finally {
      setLoading(false)
    }
  }, [collectionAddress, collectionTokenId])

  React.useEffect(() => {
    void load()
  }, [load])

  return (
    <Card>
      <CardHeader>
        <CardTitle>Agent card</CardTitle>
        <CardDescription>
          The live 6022 card served by the agent's node.
        </CardDescription>
      </CardHeader>
      <CardContent>
        {loading ? (
          <div className="space-y-3">
            <Skeleton className="h-4 w-40" />
            <Skeleton className="h-4 w-64" />
            <Skeleton className="h-4 w-52" />
          </div>
        ) : unreachable || !card ? (
          <EmptyState
            icon={Unplug}
            title="Agent node unreachable"
            description="The agent's node did not answer with its well-known card."
            action={
              <Button variant="outline" size="sm" onClick={load}>
                Retry
              </Button>
            }
          />
        ) : (
          <div className="flex flex-col gap-4">
            <CardField
              label="Protocol"
              value={`${card.protocol} ${card.protocolVersion}`}
            />
            <CardField label="Name" value={card.name} />
            <CardField label="ENS domain" value={card.ensDomain} />
            <CardField
              label="Responses endpoint"
              value={card.endpoints?.responses}
            />
            <CardField
              label="Chat completions endpoint"
              value={card.endpoints?.chatCompletions}
            />
            <div className="flex flex-col gap-1">
              <span className="text-xs text-muted-foreground">
                Payment methods
              </span>
              {card.paymentMethods && card.paymentMethods.length > 0 ? (
                <div className="flex flex-wrap gap-2">
                  {card.paymentMethods.map((method) => (
                    <Badge key={method} variant="outline">
                      {method}
                    </Badge>
                  ))}
                </div>
              ) : (
                <span className="text-sm">Free</span>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
