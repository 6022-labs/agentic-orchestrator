import React from "react";
import { useForm } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { zodResolver } from "@hookform/resolvers/zod";
import { LLMConnectionPlatformType } from "@/models/enums/llm-connection-platform-type";
import {
  SetupAgentLLMSettingsFormSchemaType,
  createSetupAgentLLMSettingsFormSchema,
} from "./setup-agent-llm-settings-form-schema";
import {
  LLMConnectionPlatform,
  LLMConnectionPlatformDisplayNames,
  getAllCloudLLMConnectionPlatforms,
  getAllSelfHostedLLMConnectionPlatforms,
} from "@/models/enums/llm-connection-platform";
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
  FormMessage,
  FormControl,
  FormDescription,
} from "@/components/ui/form";
import type { FormCompletionResult } from "@/models/api/validation-error";

type SetupAgentLLMSettingsFormProps = {
  onFormComplete: (
    values: SetupAgentLLMSettingsFormSchemaType
  ) => Promise<FormCompletionResult<keyof SetupAgentLLMSettingsFormSchemaType>>;
  defaultValues?: SetupAgentLLMSettingsFormSchemaType;
};

export default function SetupAgentLLMSettingsForm(
  props: Readonly<SetupAgentLLMSettingsFormProps>
) {
  const setupAgentLLMSettingsFormSchema =
    createSetupAgentLLMSettingsFormSchema();

  const form = useForm<SetupAgentLLMSettingsFormSchemaType>({
    resolver: zodResolver(setupAgentLLMSettingsFormSchema),
    defaultValues: props.defaultValues,
  });

  const [buttonLoading, setButtonLoading] = React.useState(false);

  const selectedLLMConnectionPlatform = form.watch("platform");
  const selectedLLMConnectionPlatformType = form.watch("platformType");

  const [availableLLMConnectionPlatforms, setAvailableLLMConnectionPlatforms] =
    React.useState<LLMConnectionPlatform[]>([]);

  const updateLLMConnectionPlatforms = () => {
    if (selectedLLMConnectionPlatformType === LLMConnectionPlatformType.Cloud) {
      setAvailableLLMConnectionPlatforms(getAllCloudLLMConnectionPlatforms());
      return;
    }

    if (
      selectedLLMConnectionPlatformType === LLMConnectionPlatformType.SelfHost
    ) {
      setAvailableLLMConnectionPlatforms(
        getAllSelfHostedLLMConnectionPlatforms()
      );
    }
  };

  React.useEffect(() => {
    if (props.defaultValues) {
      form.setValue("model", props.defaultValues.model);
      form.setValue("baseUrl", props.defaultValues.baseUrl);
      form.setValue(
        "platform",
        props.defaultValues.platform as LLMConnectionPlatform
      );
      form.setValue(
        "platformType",
        props.defaultValues.platformType as LLMConnectionPlatformType
      );
      form.setValue(
        "amazonPlatformRegion",
        props.defaultValues.amazonPlatformRegion
      );
    }
  }, [form, props.defaultValues]);

  React.useEffect(() => {
    updateLLMConnectionPlatforms();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedLLMConnectionPlatformType]);

  return (
    <Form {...form}>
      <form className="h-fit">
        <div className="space-y-4 px-1">
          <FormField
            control={form.control}
            name="platformType"
            render={({ field }) => (
              <FormItem>
                <FormLabel>LLM Platform type</FormLabel>
                <Select value={field.value} onValueChange={field.onChange}>
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select a llm platform type" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="cloud">Cloud</SelectItem>
                    <SelectItem value="self-hosted">Self-Hosted</SelectItem>
                  </SelectContent>
                </Select>
                <FormDescription>
                  Platform to use for the agent.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          {selectedLLMConnectionPlatformType && (
            <>
              <FormField
                control={form.control}
                name="platform"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>LLM Platform</FormLabel>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger className="w-full">
                          <SelectValue placeholder="Select a llm platform" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {availableLLMConnectionPlatforms.map((llmPlatform) => (
                          <SelectItem key={llmPlatform} value={llmPlatform}>
                            {LLMConnectionPlatformDisplayNames[llmPlatform]}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormDescription>
                      Platform to use for the agent.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </>
          )}
          {selectedLLMConnectionPlatform ===
            LLMConnectionPlatform.AmazonBedrock && (
            <FormField
              control={form.control}
              name="amazonPlatformRegion"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Amazon Platform Region</FormLabel>
                  <FormControl>
                    <Input placeholder="us-west-2" {...field} />
                  </FormControl>
                  <FormDescription>
                    Region for the Amazon Bedrock platform.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}
          {selectedLLMConnectionPlatformType ===
            LLMConnectionPlatformType.SelfHost && (
            <>
              <FormField
                control={form.control}
                name="baseUrl"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Server Base URL</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="https://my-ollama-server.com:11434"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Base URL of the LLM server.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </>
          )}
          {selectedLLMConnectionPlatformType ===
            LLMConnectionPlatformType.Cloud && (
            <FormField
              control={form.control}
              name="apiKey"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Platform API Key</FormLabel>
                  <FormControl>
                    <Input placeholder="sk-...8I0A" {...field} />
                  </FormControl>
                  <FormDescription>
                    Api key of{" "}
                    {
                      LLMConnectionPlatformDisplayNames[
                        form.watch("platform") as LLMConnectionPlatform
                      ]
                    }
                    .
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}
          <FormField
            control={form.control}
            name="model"
            render={({ field }) => (
              <FormItem>
                <FormLabel>LLM Model</FormLabel>
                <FormControl>
                  <Input placeholder="llama3.3" {...field} />
                </FormControl>
                <FormDescription>LLM model for your agent.</FormDescription>
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
                        .field as keyof SetupAgentLLMSettingsFormSchemaType,
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
