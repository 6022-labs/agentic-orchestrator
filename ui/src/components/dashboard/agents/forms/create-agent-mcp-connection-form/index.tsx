import React from "react";
import { useForm } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  MCPAuthenticationMethod,
  createAgentMCPConnectionFormSchema,
  MCPAuthenticationMethodDisplayNames,
  CreateAgentMCPConnectionFormSchemaType,
} from "./create-agent-mcp-connection-form-schema";
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
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { FormCompletionResult } from "@/models/api/validation-error";

type CreateAgentMCPConnectionFormProps = {
  onFormComplete: (
    values: CreateAgentMCPConnectionFormSchemaType
  ) => Promise<
    FormCompletionResult<keyof CreateAgentMCPConnectionFormSchemaType>
  >;
};

export default function CreateAgentMCPConnectionForm({
  onFormComplete,
}: Readonly<CreateAgentMCPConnectionFormProps>) {
  const [buttonLoading, setButtonLoading] = React.useState(false);

  const form = useForm<CreateAgentMCPConnectionFormSchemaType>({
    resolver: zodResolver(createAgentMCPConnectionFormSchema),
    defaultValues: {
      name: "",
      url: "",
      authenticationMethod: MCPAuthenticationMethod.NONE,
      headerAuthKey: undefined,
      headerAuthValue: undefined,
    },
  });

  const authenticationMethod = form.watch("authenticationMethod");

  return (
    <Form {...form}>
      <form className="h-fit">
        <div className="space-y-4 px-1">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input
                    placeholder="Enter name..."
                    {...field}
                    onChange={(e) => {
                      let value = e.target.value;
                      value = value.replace(/ /g, "-").toLowerCase();
                      value = value.replace(/[^a-z0-9_-]/g, "");

                      if (value === "" || /^[a-z][a-z0-9_-]*$/.test(value)) {
                        e.target.value = value;
                        field.onChange(e);
                      }
                    }}
                  />
                </FormControl>
                <FormDescription>Name of the MCP connection.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="url"
            render={({ field }) => (
              <FormItem>
                <FormLabel>URL</FormLabel>
                <FormControl>
                  <Input placeholder="https://my.mcp.com/v1/..." {...field} />
                </FormControl>
                <FormDescription>URL of the MCP server.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="authenticationMethod"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Authentication</FormLabel>
                <Select
                  value={field.value}
                  onValueChange={(value) => {
                    field.onChange(value);
                    form.clearErrors(["headerAuthKey", "headerAuthValue"]);
                    form.setValue("headerAuthValue", undefined);
                    if (value === MCPAuthenticationMethod.API_KEY) {
                      form.setValue("headerAuthKey", "Authorization");
                      return;
                    }

                    form.setValue("headerAuthKey", undefined);
                  }}
                >
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {Object.values(MCPAuthenticationMethod).map((method) => (
                      <SelectItem key={method} value={method}>
                        {MCPAuthenticationMethodDisplayNames[method]}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>
                  Authentication method for the MCP server.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          {authenticationMethod === MCPAuthenticationMethod.CUSTOM_HEADER && (
            <FormField
              control={form.control}
              name="headerAuthKey"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Header key</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter header key..." {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}
          {authenticationMethod !== MCPAuthenticationMethod.NONE && (
            <FormField
              control={form.control}
              name="headerAuthValue"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {authenticationMethod === MCPAuthenticationMethod.API_KEY
                      ? "API Key"
                      : "Header value"}
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={
                        authenticationMethod === MCPAuthenticationMethod.API_KEY
                          ? "Enter API key..."
                          : "Enter header value..."
                      }
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}
          <div className="flex w-full justify-end gap-x-2">
            <Button
              type="button"
              loading={buttonLoading}
              onClick={async () => {
                setButtonLoading(true);
                try {
                  form.clearErrors();

                  const valid = await form.trigger();
                  if (!valid) {
                    return;
                  }

                  const values = form.getValues();
                  if (
                    values.authenticationMethod ===
                    MCPAuthenticationMethod.API_KEY
                  ) {
                    values.headerAuthKey = "Authorization";
                    values.headerAuthValue = `Bearer ${values.headerAuthValue}`;
                  }

                  const submission = await onFormComplete(values);

                  if (submission?.error) {
                    form.setError(
                      submission.error
                        .field as keyof CreateAgentMCPConnectionFormSchemaType,
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
