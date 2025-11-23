import React from "react";
import { cn } from "@/lib/utils";
import { CircleAlert } from "lucide-react";
import { Button } from "@/components/ui/button";
import { CopyBlock } from "@/components/ui/copy-block";
import MySQLSetupScript from "./setup-scripts/mysql-setup-script";
import MariaDBSetupScript from "./setup-scripts/mariadb-setup-script";
import PostgresSetupScript from "./setup-scripts/postgres-setup-script";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { DatabaseConnectionDriver } from "@/models/enums/database-connection-driver";
import {
  FormControl,
  FormDescription,
  FormItem,
  FormLabel,
} from "@/components/ui/form";

export type SetupAgentDatabaseFirstStepResult = {
  databaseName: string;
  databaseUsername: string;
  databasePassword: string;
};

type SetupAgentDatabaseFirstStepResultInformationProps = {
  className?: string;
  driver: DatabaseConnectionDriver;
  result?: SetupAgentDatabaseFirstStepResult;
  onNextValidStep: () => Promise<void>;
};

export function SetupAgentDatabaseFirstStepResultInformation(
  props: Readonly<SetupAgentDatabaseFirstStepResultInformationProps>
) {
  const className = cn("flex flex-col space-y-4 px-1", props.className);

  const [numberOfClicksOnValidate, setNumberOfClicksOnValidate] =
    React.useState<number>(0);

  const warningRef = React.useRef<HTMLDivElement>(null);

  const setupScripts = {
    [DatabaseConnectionDriver.None]: <></>,
    [DatabaseConnectionDriver.Postgres]: (
      <PostgresSetupScript result={props.result} />
    ),
    [DatabaseConnectionDriver.MySQL]: <MySQLSetupScript />,
    [DatabaseConnectionDriver.MariaDB]: <MariaDBSetupScript />,
  };

  return (
    <div className={className}>
      <div className="flex flex-col space-y-4">
        <FormItem>
          <FormLabel>Database name</FormLabel>
          <FormControl>
            <CopyBlock>{props.result?.databaseName ?? ""}</CopyBlock>
          </FormControl>
          <FormDescription>
            This is the database you must create.
          </FormDescription>
        </FormItem>
        <FormItem>
          <FormLabel>Database username</FormLabel>
          <FormControl>
            <CopyBlock>{props.result?.databaseUsername ?? ""}</CopyBlock>
          </FormControl>
          <FormDescription>
            Role you must create in your database.
          </FormDescription>
        </FormItem>
        <FormItem>
          <FormLabel>Database password</FormLabel>
          <FormControl>
            <CopyBlock>{props.result?.databasePassword ?? ""}</CopyBlock>
          </FormControl>
          <FormDescription>
            Password you must use for the role in your database.
          </FormDescription>
        </FormItem>
      </div>
      <div className="relative overflow-hidden">
        {setupScripts[props.driver]}
      </div>
      <Alert ref={warningRef}>
        <CircleAlert className="h-4 w-4" />
        <AlertTitle>Caution!</AlertTitle>
        <AlertDescription>
          Those credentials will not be shown again. Make sure to back them up.
        </AlertDescription>
      </Alert>
      <div className="flex w-full justify-end gap-x-2">
        <Button
          type="button"
          onClick={async (e) => {
            if (numberOfClicksOnValidate == 0) {
              setNumberOfClicksOnValidate(1);
              e.currentTarget.classList.add("animate-shake");
              warningRef.current?.classList.add("animate-shake");
            } else {
              await props.onNextValidStep();
            }
          }}
        >
          {numberOfClicksOnValidate == 0
            ? "I've backed up credentials"
            : "Yes, I really have"}
        </Button>
      </div>
    </div>
  );
}
