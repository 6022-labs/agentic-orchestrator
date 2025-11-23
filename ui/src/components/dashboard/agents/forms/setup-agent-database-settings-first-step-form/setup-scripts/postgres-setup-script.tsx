import { Label } from "@/components/ui/label";
import { CopyButton } from "@/components/ui/copy-button";
import { SetupAgentDatabaseFirstStepResult } from "../setup-agent-database-first-step-result-information";
import {
  CodeBlock,
  CodeBlockBody,
  CodeBlockHeader,
  CodeBlockContent,
  CodeBlockFilename,
} from "@/components/ui/code-block";

type PostgresSetupScriptProps = {
  className?: string;
  result?: SetupAgentDatabaseFirstStepResult;
};

const firstBlockCodeTemplate = `-- Connect to the default database to create the new database

-- Run this command separately to avoid transaction error
CREATE DATABASE "<DATABASE_NAME>";`;

const secondBlockCodeTemplate = `-- Then connect to the newly created database and run the following command

-- Create the role with login privilege
CREATE ROLE "<ROLE_USERNAME>"
WITH LOGIN PASSWORD '<ROLE_PASSWORD>';

-- Grant all privileges on the database to the role
GRANT ALL PRIVILEGES
ON DATABASE "<DATABASE_NAME>"
TO "<ROLE_USERNAME>";

GRANT ALL PRIVILEGES
ON SCHEMA public
TO "<ROLE_USERNAME>";`;

export default function PostgresSetupScript(props: PostgresSetupScriptProps) {
  const firstBlockCode = firstBlockCodeTemplate.replace(
    /<DATABASE_NAME>/g,
    props.result?.databaseName ?? ""
  );

  const secondBlockCode = secondBlockCodeTemplate
    .replace(/<DATABASE_NAME>/g, props.result?.databaseName ?? "")
    .replace(/<ROLE_USERNAME>/g, props.result?.databaseUsername ?? "")
    .replace(/<ROLE_PASSWORD>/g, props.result?.databasePassword ?? "");

  return (
    <div className="grid grid-rows-1 space-y-4">
      <Label>
        Postgres setup script
        <span className="text-xs text-muted-foreground">
          (Run these commands in your Postgres client)
        </span>
      </Label>
      <CodeBlock defaultValue="sql">
        <CodeBlockHeader>
          <CodeBlockFilename value="sql">
            postgres-create-database-script.sql
          </CodeBlockFilename>
          <CopyButton text={firstBlockCode} />
        </CodeBlockHeader>
        <CodeBlockBody value="sql">
          <CodeBlockContent language="sql">{firstBlockCode}</CodeBlockContent>
        </CodeBlockBody>
      </CodeBlock>
      <CodeBlock defaultValue="sql">
        <CodeBlockHeader>
          <CodeBlockFilename value="sql">
            postgres-create-role-script.sql
          </CodeBlockFilename>
          <CopyButton text={secondBlockCode} />
        </CodeBlockHeader>
        <CodeBlockBody value="sql">
          <CodeBlockContent language="sql">{secondBlockCode}</CodeBlockContent>
        </CodeBlockBody>
      </CodeBlock>
    </div>
  );
}
