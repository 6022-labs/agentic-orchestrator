import React from "react"
import axios from "axios"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Skeleton } from "@/components/ui/skeleton"
import { Textarea } from "@/components/ui/textarea"
import { extractErrorMessage } from "@/lib/api-errors"
import type { AgentRole } from "@/models/enums/agent-role"
import { PROMPT_TEMPLATE_KEYS } from "@/models/prompt-template"
import {
  getPromptTemplate,
  savePromptTemplate,
} from "@/queries/api/prompt-templates.queries"

// Keys matching the agent's role come first so the relevant templates are on top.
const orderedKeys = (role: AgentRole) =>
  [...PROMPT_TEMPLATE_KEYS].sort((a, b) => {
    const aMatches = a.startsWith(role) ? 0 : 1
    const bMatches = b.startsWith(role) ? 0 : 1
    return aMatches - bMatches
  })

export default function PromptTemplatesPanel({
  collectionAddress,
  collectionTokenId,
  role,
}: {
  collectionAddress: string
  collectionTokenId: string
  role: AgentRole
}) {
  const keys = React.useMemo(() => orderedKeys(role), [role])
  const [selectedKey, setSelectedKey] = React.useState<string>(keys[0])
  const [content, setContent] = React.useState("")
  const [loading, setLoading] = React.useState(true)
  const [saving, setSaving] = React.useState(false)

  const load = React.useCallback(async () => {
    setLoading(true)
    try {
      const { data } = await getPromptTemplate(
        collectionAddress,
        collectionTokenId,
        selectedKey
      )
      setContent(data.content ?? "")
    } catch (error) {
      // 404 means the template is unset — start from an empty editor.
      if (axios.isAxiosError(error) && error.response?.status === 404) {
        setContent("")
      } else {
        toast.error(extractErrorMessage(error, "Failed to load the template"))
      }
    } finally {
      setLoading(false)
    }
  }, [collectionAddress, collectionTokenId, selectedKey])

  React.useEffect(() => {
    void load()
  }, [load])

  const onSave = async () => {
    setSaving(true)
    try {
      await savePromptTemplate(
        collectionAddress,
        collectionTokenId,
        selectedKey,
        content
      )
      toast.success("Template saved")
    } catch (error) {
      toast.error(extractErrorMessage(error, "Failed to save the template"))
    } finally {
      setSaving(false)
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Prompt templates</CardTitle>
        <CardDescription>
          Per-agent overrides of the orchestrator's default prompts.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        <Select value={selectedKey} onValueChange={setSelectedKey}>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Template key" />
          </SelectTrigger>
          <SelectContent>
            {keys.map((key) => (
              <SelectItem key={key} value={key}>
                {key}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        {loading ? (
          <Skeleton className="h-40 w-full" />
        ) : (
          <Textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="No template set — the orchestrator default applies."
            className="min-h-40"
          />
        )}

        <div className="flex justify-end">
          <Button onClick={onSave} loading={saving} disabled={loading}>
            Save template
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
