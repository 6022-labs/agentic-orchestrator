import React from "react";
import { useForm } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  CreateAgentContextFormSchemaType,
  createCreateAgentContextFormSchema,
} from "./create-agent-context-form-schema";
import type { FormCompletionResult } from "@/models/api/validation-error";
import {
  Form,
  FormItem,
  FormLabel,
  FormField,
  FormMessage,
  FormControl,
  FormDescription,
} from "@/components/ui/form";

type CreateAgentContextFormProps = {
  onFormComplete: (
    values: CreateAgentContextFormSchemaType
  ) => Promise<FormCompletionResult<keyof CreateAgentContextFormSchemaType>>;
};

export default function CreateAgentContextForm(
  props: Readonly<CreateAgentContextFormProps>
) {
  const [buttonLoading, setButtonLoading] = React.useState(false);
  const createAgentContextFormSchema = createCreateAgentContextFormSchema();

  const form = useForm<CreateAgentContextFormSchemaType>({
    resolver: zodResolver(createAgentContextFormSchema),
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
                  <Input placeholder="Enter key..." {...field} />
                </FormControl>
                <FormDescription>
                  The key is used to identify the context. It should be unique.
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
                  const keyIsValidPromise = form.trigger("key");

                  const errors = await Promise.all([keyIsValidPromise]);

                  if (errors.some((error) => error === false)) {
                    return;
                  }

                  const submission = await props.onFormComplete(
                    form.getValues()
                  );

                  if (submission?.error) {
                    form.setError(
                      submission.error
                        .field as keyof CreateAgentContextFormSchemaType,
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
              Create
            </Button>
          </div>
        </div>
      </form>
    </Form>
  );
}
