import React from "react";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  AgentContextMessageRole,
  AgentContextMessageRoleDisplayNames,
} from "@/models/enums/agent-context-message-role";
import {
  AddMessageToAgentContextFormSchemaType,
  createAddMessageToAgentContextFormSchema,
} from "./add-message-to-agent-context-form-schema";
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
import {
  Select,
  SelectItem,
  SelectValue,
  SelectContent,
  SelectTrigger,
} from "@/components/ui/select";

type AddMessageToAgentContextFormProps = {
  onFormComplete: (
    values: AddMessageToAgentContextFormSchemaType
  ) => Promise<
    FormCompletionResult<keyof AddMessageToAgentContextFormSchemaType>
  >;
};

export default function AddMessageToAgentContextForm(
  props: Readonly<AddMessageToAgentContextFormProps>
) {
  const [buttonLoading, setButtonLoading] = React.useState(false);
  const addMessageToAgentContextFormSchema =
    createAddMessageToAgentContextFormSchema();

  const form = useForm<AddMessageToAgentContextFormSchemaType>({
    resolver: zodResolver(addMessageToAgentContextFormSchema),
  });

  return (
    <Form {...form}>
      <form className="h-fit">
        <div className="space-y-4 px-1">
          <FormField
            control={form.control}
            name="role"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Role</FormLabel>
                <Select value={field.value} onValueChange={field.onChange}>
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select a role" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {Object.values(AgentContextMessageRole).map((role) => (
                      <SelectItem key={role} value={role}>
                        {AgentContextMessageRoleDisplayNames[role]}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>
                  Role of the message to add to the context.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="message"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Message</FormLabel>
                <FormControl>
                  <Textarea
                    placeholder="I've take in account the following..."
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  Message to add in the context.
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
                  const roleIsValidPromise = form.trigger("role");
                  const messageIsValidPromise = form.trigger("message");

                  const errors = await Promise.all([
                    roleIsValidPromise,
                    messageIsValidPromise,
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
                        .field as keyof AddMessageToAgentContextFormSchemaType,
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
