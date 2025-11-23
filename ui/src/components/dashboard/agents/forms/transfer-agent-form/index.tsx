import React from "react";
import { useForm } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  TransferAgentFormSchemaType,
  createTransferAgentFormSchema,
} from "./transfer-agent-form-schema";
import {
  Form,
  FormItem,
  FormLabel,
  FormField,
  FormMessage,
  FormControl,
  FormDescription,
} from "@/components/ui/form";
import type { FormCompletionResult } from "@/models/api/validation-error";

type TransferAgentFormProps = {
  onFormComplete: (
    values: TransferAgentFormSchemaType
  ) => Promise<FormCompletionResult<keyof TransferAgentFormSchemaType>>;
};

export default function TransferAgentForm(
  props: Readonly<TransferAgentFormProps>
) {
  const transferAgentFormSchema = createTransferAgentFormSchema();

  const form = useForm<TransferAgentFormSchemaType>({
    resolver: zodResolver(transferAgentFormSchema),
  });

  const [buttonLoading, setButtonLoading] = React.useState(false);

  return (
    <Form {...form}>
      <form className="h-fit">
        <div className="space-y-4 px-1">
          <FormField
            control={form.control}
            name="address"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Receiver</FormLabel>
                <FormControl>
                  <Input placeholder="0x..." {...field} />
                </FormControl>
                <FormDescription>Address of the receiver.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <div className="flex w-full justify-end gap-x-2">
            <Button
              type="button"
              loading={buttonLoading}
              disabled={buttonLoading}
              onClick={async () => {
                setButtonLoading(true);

                try {
                  const addressIsValidPromise = form.trigger("address");

                  const errors = await Promise.all([addressIsValidPromise]);

                  if (errors.some((error) => error === false)) {
                    return;
                  }

                  const submission = await props.onFormComplete(
                    form.getValues()
                  );

                  if (submission?.error) {
                    form.setError(
                      submission.error
                        .field as keyof TransferAgentFormSchemaType,
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
              Transfer
            </Button>
          </div>
        </div>
      </form>
    </Form>
  );
}
