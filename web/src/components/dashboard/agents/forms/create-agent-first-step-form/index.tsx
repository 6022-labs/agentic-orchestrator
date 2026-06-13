import React from "react";
import { Eip1193Provider } from "ethers";
import { useForm } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAppKitProvider } from "@reown/appkit/react";
import { useAgentCollections } from "@/contexts/agent-collections-context";
import type { FormCompletionResult } from "@/models/api/validation-error";
import {
  CreateAgentFirstStepDisabledValues,
  CreateAgentFirstStepFormSchemaType,
  createCreateAgentFirstStepFormSchema,
} from "./create-agent-first-step-form-schema";
import {
  AgentRole,
  getAllAgentRoles,
  AgentRoleDisplayNames,
} from "@/models/enums/agent-role";
import {
  Select,
  SelectItem,
  SelectValue,
  SelectContent,
  SelectTrigger,
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

type CreateAgentFirstStepFormProps = {
  accessibleRoles?: AgentRole[];
  disabledValues?: CreateAgentFirstStepDisabledValues;
  defaultValues?: Partial<CreateAgentFirstStepFormSchemaType>;
  onFormComplete: (
    values: CreateAgentFirstStepFormSchemaType
  ) => Promise<FormCompletionResult<keyof CreateAgentFirstStepFormSchemaType>>;
};

export default function CreateAgentFirstStepForm({
  onFormComplete,
  defaultValues = {},
  disabledValues = {},
  accessibleRoles = getAllAgentRoles(),
}: Readonly<CreateAgentFirstStepFormProps>) {
  const [buttonLoading, setButtonLoading] = React.useState(false);

  const { walletProvider } = useAppKitProvider("eip155");
  const eip1193Provider = walletProvider as Eip1193Provider;

  const createAgentFirstStepFormSchema =
    createCreateAgentFirstStepFormSchema(eip1193Provider);

  const form = useForm<CreateAgentFirstStepFormSchemaType>({
    resolver: zodResolver(createAgentFirstStepFormSchema),
    defaultValues: defaultValues as CreateAgentFirstStepFormSchemaType,
  });

  const { collections } = useAgentCollections();

  return (
    <Form {...form}>
      <form className="h-fit">
        <div className="space-y-4 px-1">
          <FormField
            control={form.control}
            name="collectionAddress"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Collection</FormLabel>
                <Select
                  value={field.value}
                  onValueChange={field.onChange}
                  disabled={disabledValues.collectionAddress}
                >
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue
                        defaultValue={defaultValues?.collectionAddress}
                        placeholder="Select a collection"
                      />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {collections.map((collection) => (
                      <SelectItem
                        key={collection.address}
                        value={collection.address}
                      >
                        {collection.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>Collection of your agent.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
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
                    disabled={disabledValues.name}
                    defaultValue={defaultValues?.name}
                    onChange={(e) => {
                      e.target.value = e.target.value.toLowerCase();
                      e.target.value = e.target.value.replace(/ /g, "-");

                      const regex = /^([a-z0-9]|-)*$/gm;
                      if (regex.test(e.target.value)) field.onChange(e);
                    }}
                  />
                </FormControl>
                <FormDescription>This is the agent name.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="role"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Role</FormLabel>
                <Select
                  value={field.value}
                  onValueChange={field.onChange}
                  disabled={disabledValues.role}
                >
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue
                        placeholder="Select a role"
                        defaultValue={defaultValues?.role ?? AgentRole.None}
                      />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {Object.values(accessibleRoles).map((role) => (
                      <SelectItem key={role} value={role}>
                        {AgentRoleDisplayNames[role]}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>Role of your agent.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          {form.watch("role") === AgentRole.Human && (
            <FormField
              control={form.control}
              name="walletAddress"
              render={({ field }) => (
                <FormItem>
                  <FormLabel htmlFor="walletAddress">Wallet Address</FormLabel>
                  <FormControl>
                    <Input
                      id="walletAddress"
                      placeholder="0x1234...abcd"
                      {...field}
                      disabled={disabledValues.walletAddress}
                      defaultValue={defaultValues?.walletAddress}
                    />
                  </FormControl>
                  <FormDescription>
                    Wallet address of the human agent.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}
          {form.watch("role") === AgentRole.Clone && (
            <FormField
              control={form.control}
              name="cloneOfAgentId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel htmlFor="cloneOfAgentId">
                    Clone of Agent ID
                  </FormLabel>
                  <FormControl>
                    <Input
                      id="cloneOfAgentId"
                      placeholder="1"
                      {...field}
                      disabled={disabledValues.cloneOfAgentId}
                      defaultValue={defaultValues?.cloneOfAgentId}
                    />
                  </FormControl>
                  <FormDescription>
                    The ID of the agent to clone.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}
          <FormField
            control={form.control}
            name="mission"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Mission</FormLabel>
                <FormControl>
                  <Textarea
                    disabled={disabledValues.mission}
                    defaultValue={defaultValues?.mission}
                    placeholder="Your mission is to help users with their tasks..."
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  Mission of your agent. This will be used to generate the
                  agent's response.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="image"
            render={({ field }) => (
              <FormItem>
                <FormLabel htmlFor="image">Image</FormLabel>
                <FormControl>
                  <Input
                    id="image"
                    type="file"
                    ref={field.ref}
                    name={field.name}
                    onBlur={field.onBlur}
                    disabled={disabledValues.image}
                    accept="image/png, image/gif, image/jpg, image/jpeg"
                    onChange={(e) => field.onChange(e.target.files?.[0])}
                  />
                </FormControl>
                <FormDescription>This is the agent image.</FormDescription>
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
                  const isValid = await form.trigger();
                  if (!isValid) {
                    return;
                  }

                  const values = form.getValues();
                  const submission = await onFormComplete(values);

                  if (submission?.error) {
                    form.setError(
                      submission.error
                        .field as keyof CreateAgentFirstStepFormSchemaType,
                      {
                        type: "server",
                        message: submission.error.message,
                      }
                    );
                    return;
                  }
                } catch (error) {
                  console.error(error);
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
