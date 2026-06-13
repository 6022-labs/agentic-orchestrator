import React from "react";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  UpdateAgentPromptFormSchemaType,
  createUpdateAgentPromptFormSchema,
} from "./update-agent-prompt-form-schema";
import {
  AgentPromptTemplateKey,
  AgentPromptTemplateKeyDisplayNames,
} from "@/models/enums/agent-prompt-template-key";
import {
  Select,
  SelectItem,
  SelectValue,
  SelectTrigger,
  SelectContent,
} from "@/components/ui/select";
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

type UpdateAgentPromptFormProps = {
  onFormComplete: (
    values: UpdateAgentPromptFormSchemaType
  ) => Promise<FormCompletionResult<keyof UpdateAgentPromptFormSchemaType>>;
  keys: AgentPromptTemplateKey[];
  defaultPromptsTemplates?: Map<string, string>;
};

export default function UpdateAgentPromptForm(
  props: Readonly<UpdateAgentPromptFormProps>
) {
  const [buttonLoading, setButtonLoading] = React.useState(false);

  const updateAgentPromptFormSchema = createUpdateAgentPromptFormSchema();

  const form = useForm<UpdateAgentPromptFormSchemaType>({
    resolver: zodResolver(updateAgentPromptFormSchema),
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
                <FormLabel>Prompt type</FormLabel>
                <Select
                  value={field.value}
                  onValueChange={(value) => {
                    form.setValue(
                      "content",
                      props.defaultPromptsTemplates?.get(value) ?? ""
                    );
                    field.onChange(value);
                  }}
                >
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select a prompt type" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {props.keys.map((key) => (
                      <SelectItem key={key} value={key}>
                        {AgentPromptTemplateKeyDisplayNames[key]}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>Prompt to update.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="content"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Prompt</FormLabel>
                <FormControl>
                  <Textarea
                    placeholder="Your a an AI assistant..."
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  Prompt to use for the agent. This will be used to generate the
                  agent's response.
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
                  const contentIsValidPromise = form.trigger("content");

                  const errors = await Promise.all([
                    keyIsValidPromise,
                    contentIsValidPromise,
                  ]);

                  if (errors.some((error) => error === false)) {
                    return;
                  }

                  const submission = await props.onFormComplete(
                    form.getValues()
                  );

                  if (submission?.error) {
                    form.setError(
                      submission.error
                        .field as keyof UpdateAgentPromptFormSchemaType,
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
