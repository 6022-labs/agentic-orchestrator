import React from "react";
import { cn } from "@/lib/utils";
import { Input } from "@/components/ui/input";
import { UseFormReturn } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { SetupAgentDatabaseSettingsFirstStepFormSchemaType } from "./setup-agent-database-settings-first-step-form-schema";
import {
  getAllDatabaseConnectionDrivers,
  DatabaseConnectionDriverDisplayNames,
} from "@/models/enums/database-connection-driver";
import {
  FormItem,
  FormLabel,
  FormField,
  FormControl,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import {
  Select,
  SelectItem,
  SelectValue,
  SelectContent,
  SelectTrigger,
} from "@/components/ui/select";

type SetupAgentDatabaseFirstStepInformationFormProps = {
  form: UseFormReturn<SetupAgentDatabaseSettingsFirstStepFormSchemaType>;
  disabledValues?: Map<string, boolean>;

  nextButtonLabel: string;
  onNextValidStep: () => Promise<void>;

  className?: string;
};

export function SetupAgentDatabaseFirstStepInformationForm(
  props: Readonly<SetupAgentDatabaseFirstStepInformationFormProps>
) {
  const [buttonLoading, setButtonLoading] = React.useState(false);

  return (
    <div className={cn(props.className, "space-y-4 px-1")}>
      <FormField
        control={props.form.control}
        name="host"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Database Host</FormLabel>
            <FormControl>
              <Input
                disabled={props.disabledValues?.get("host")}
                placeholder="my-database-server.com"
                {...field}
              />
            </FormControl>
            <FormDescription>Host of the database server.</FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={props.form.control}
        name="port"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Database Port</FormLabel>
            <FormControl>
              <Input
                disabled={props.disabledValues?.get("port")}
                placeholder="6689"
                type="number"
                {...field}
                onChange={(e) => {
                  const value = e.target.value;

                  if (/^\d*$/.test(value)) {
                    field.onChange(value ? parseInt(value, 10) : "");
                  }
                }}
              />
            </FormControl>
            <FormDescription>Port of the database server.</FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={props.form.control}
        name="driver"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Database driver</FormLabel>
            <Select
              value={field.value}
              onValueChange={field.onChange}
              disabled={props.disabledValues?.get("driver")}
            >
              <FormControl>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Select a database driver" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {getAllDatabaseConnectionDrivers().map((databaseDriver) => (
                  <SelectItem key={databaseDriver} value={databaseDriver}>
                    {DatabaseConnectionDriverDisplayNames[databaseDriver]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormDescription>Database driver of the agent.</FormDescription>
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
              const hostIsValidPromise = props.form.trigger("host");
              const portIsValidPromise = props.form.trigger("port");
              const driverIsValidPromise = props.form.trigger("driver");

              const errors = await Promise.all([
                hostIsValidPromise,
                portIsValidPromise,
                driverIsValidPromise,
              ]);

              if (errors.some((error) => error === false)) {
                return;
              }

              await props.onNextValidStep();
            } finally {
              setButtonLoading(false);
            }
          }}
        >
          {props.nextButtonLabel}
        </Button>
      </div>
    </div>
  );
}
