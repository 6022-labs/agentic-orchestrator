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
import type { FormCompletionResult } from "@/models/api/validation-error"
import {
  transferAgentFormSchema,
  type TransferAgentFormSchemaType,
} from "./transfer-agent-form-schema"

type Props = {
  onFormComplete: (
    values: TransferAgentFormSchemaType
  ) => Promise<FormCompletionResult<keyof TransferAgentFormSchemaType>>
}

export default function TransferAgentForm({ onFormComplete }: Props) {
  const form = useForm<TransferAgentFormSchemaType>({
    resolver: zodResolver(transferAgentFormSchema),
    defaultValues: { to: "" },
  })

  const onSubmit = async (values: TransferAgentFormSchemaType) => {
    const { error } = await onFormComplete(values)

    if (error) {
      form.setError(error.field, { type: "server", message: error.message })
    }
  }

  return (
    <Form {...form}>
      <form
        id="transfer-agent-form"
        onSubmit={form.handleSubmit(onSubmit)}
        className="space-y-4"
      >
        <FormField
          control={form.control}
          name="to"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Recipient</FormLabel>
              <FormControl>
                <Input placeholder="0x…" {...field} />
              </FormControl>
              <FormDescription>
                The wallet that becomes the agent's owner.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      </form>
    </Form>
  )
}
