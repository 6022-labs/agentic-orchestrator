import React from "react";
import { Eip1193Provider } from "ethers";
import { useForm } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAppKitProvider } from "@reown/appkit/react";
import {
  CreateAgentCollectionFormSchemaType,
  createCreateAgentCollectionFormSchema,
} from "./create-agent-collection-form-schema";
import type { FormCompletionResult } from "@/models/api/validation-error";
import {
  Form,
  FormItem,
  FormLabel,
  FormField,
  FormControl,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";

type CreateAgentCollectionFormProps = {
  onFormComplete: (
    values: CreateAgentCollectionFormSchemaType
  ) => Promise<FormCompletionResult<keyof CreateAgentCollectionFormSchemaType>>;
};

export default function CreateAgentCollectionForm(
  props: Readonly<CreateAgentCollectionFormProps>
) {
  const [buttonLoading, setButtonLoading] = React.useState(false);

  const { walletProvider } = useAppKitProvider("eip155");
  const eip1193Provider = walletProvider as Eip1193Provider;

  const createAgentCollectionFormSchema =
    createCreateAgentCollectionFormSchema(eip1193Provider);

  const form = useForm<CreateAgentCollectionFormSchemaType>({
    resolver: zodResolver(createAgentCollectionFormSchema),
  });

  return (
    <Form {...form}>
      <form className="h-fit">
        <div className="space-y-4 px-1">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel htmlFor="name">Name</FormLabel>
                <FormControl>
                  <Input
                    id="name"
                    placeholder="John Doe"
                    {...field}
                    onChange={(e) => {
                      e.target.value = e.target.value.toLowerCase();
                      e.target.value = e.target.value.replace(/ /g, "-");

                      const regex = /^([a-z0-9]|-)*$/gm;
                      if (regex.test(e.target.value)) field.onChange(e);
                    }}
                  />
                </FormControl>
                <FormDescription>This is the collection name.</FormDescription>
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
                  <Input id="symbol" placeholder="AGT" {...field} />
                </FormControl>
                <FormDescription>
                  This is the collection symbol.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="description"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Description</FormLabel>
                <FormControl>
                  <Textarea placeholder="Your description here..." {...field} />
                </FormControl>
                <FormDescription>
                  Description of your collection. This will be displayed in
                  OpenSea, Rarible, etc.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <div className="flex w-full justify-end gap-x-2">
            <Button
              type="button"
              disabled={buttonLoading}
              onClick={async () => {
                setButtonLoading(true);

                try {
                  const nameIsValidPromise = form.trigger("name");
                  const symbolIsValidPromise = form.trigger("symbol");
                  const descriptionIsValidPromise = form.trigger("description");

                  const errors = await Promise.all([
                    nameIsValidPromise,
                    symbolIsValidPromise,
                    descriptionIsValidPromise,
                  ]);

                  if (errors.some((error) => error === false)) {
                    return;
                  }

                  const values =
                    form.getValues() as CreateAgentCollectionFormSchemaType;

                  const submission = await props.onFormComplete(values);

                  if (submission?.error) {
                    form.setError(
                      submission.error
                        .field as keyof CreateAgentCollectionFormSchemaType,
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
