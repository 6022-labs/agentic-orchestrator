import React from "react";
import { useForm } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { zodResolver } from "@hookform/resolvers/zod";
import { AutocompleteInput } from "@/components/ui/auto-complete";
import {
  UpsertAgentAttributeFormSchema,
  createUpsertAgentAttributeFormSchema,
} from "./upsert-agent-attribute-form-schema";
import {
  Form,
  FormItem,
  FormLabel,
  FormField,
  FormControl,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import type { FormCompletionResult } from "@/models/api/validation-error";

type UpsertAgentAttributeFormProps = {
  onFormComplete: (
    values: UpsertAgentAttributeFormSchema
  ) => Promise<FormCompletionResult<keyof UpsertAgentAttributeFormSchema>>;
  existingAttributes: Record<string, string>;
  defaultPromptsTemplates?: Map<string, string>;
};

export default function UpsertAgentAttributeForm(
  props: Readonly<UpsertAgentAttributeFormProps>
) {
  const [buttonLoading, setButtonLoading] = React.useState(false);

  const upsertAgentAttributeFormSchema = createUpsertAgentAttributeFormSchema();

  const form = useForm<UpsertAgentAttributeFormSchema>({
    resolver: zodResolver(upsertAgentAttributeFormSchema),
  });

  return (
    <Form {...form}>
      <form className="h-fit">
        <div className="space-y-4 px-1">
          <FormField
            control={form.control}
            name="key"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Key</FormLabel>
                <FormControl>
                  <AutocompleteInput
                    options={Object.keys(props.existingAttributes).map(
                      (key) => ({
                        label: key,
                        value: key,
                      })
                    )}
                    value={field.value}
                    onValueChange={(value) => {
                      if (value === "") {
                        field.onChange(value);
                        return;
                      }

                      const regex = /^[A-Za-z0-9_-]+$/gm;
                      if (regex.test(value)) field.onChange(value);
                    }}
                    placeholder="Select or create a key..."
                  />
                </FormControl>
                <FormDescription>Key to update.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="value"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Value</FormLabel>
                <FormControl>
                  <Input id="value" placeholder="Enter value" {...field} />
                </FormControl>
                <FormDescription>
                  Value to set for the attribute.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <div className="flex w-full justify-end gap-x-2">
            <Button
              type="button"
              loading={buttonLoading}
              onClick={async () => {
                setButtonLoading(true);

                try {
                  const isValid = await form.trigger();

                  if (!isValid) {
                    return;
                  }

                  const submission = await props.onFormComplete(
                    form.getValues()
                  );

                  if (submission?.error) {
                    form.setError(
                      submission.error
                        .field as keyof UpsertAgentAttributeFormSchema,
                      {
                        type: "server",
                        message: submission.error.message,
                      }
                    );
                    return;
                  }
                } finally {
                  setButtonLoading(false);
                }
              }}
            >
              Save
            </Button>
          </div>
        </div>
      </form>
    </Form>
  );
}
