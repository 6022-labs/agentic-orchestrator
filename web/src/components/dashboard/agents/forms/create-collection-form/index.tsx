import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"

import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { normalizeEnsLabel } from "@/lib/ens-label"
import { useAuth } from "@/contexts/auth-context"
import type { FormCompletionResult } from "@/models/api/validation-error"
import {
  createCollectionFormSchema,
  type CreateCollectionFormSchemaType,
} from "./create-collection-form-schema"

type Props = {
  onFormComplete: (
    values: CreateCollectionFormSchemaType
  ) => Promise<FormCompletionResult<keyof CreateCollectionFormSchemaType>>
}

export default function CreateCollectionForm({ onFormComplete }: Props) {
  const { user } = useAuth()

  const form = useForm<CreateCollectionFormSchemaType>({
    resolver: zodResolver(createCollectionFormSchema),
    defaultValues: {
      name: "",
      symbol: "",
      description: "",
      admin: user?.address ?? "",
      moderator: user?.address ?? "",
    },
  })

  // No reset on success: the dialog closes and unmounts the form. Values are
  // kept on failure so the user can correct and retry.
  const onSubmit = async (values: CreateCollectionFormSchemaType) => {
    const { error } = await onFormComplete(values)

    if (error) {
      form.setError(error.field, { type: "server", message: error.message })
    }
  }

  return (
    <Form {...form}>
      <form
        id="create-collection-form"
        onSubmit={form.handleSubmit(onSubmit)}
        className="space-y-4"
      >
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input
                  placeholder="my-agents"
                  {...field}
                  onChange={(e) =>
                    field.onChange(normalizeEnsLabel(e.target.value))
                  }
                />
              </FormControl>
              <FormDescription>
                Also the collection's ENS label — lowercase, digits, hyphens.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="symbol"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Symbol</FormLabel>
              <FormControl>
                <Input placeholder="AGNT" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description (optional)</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="What this collection groups…"
                  className="min-h-16"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="admin"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Admin</FormLabel>
              <FormControl>
                <Input placeholder="0x…" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="moderator"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Moderator</FormLabel>
              <FormControl>
                <Input placeholder="0x…" {...field} />
              </FormControl>
              <FormDescription>
                Moderated mints need this address's approval.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      </form>
    </Form>
  )
}
