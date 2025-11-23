import {
  CodeBlock,
  CodeBlockBody,
  CodeBlockHeader,
  CodeBlockContent,
  CodeBlockFilename,
  CodeBlockCopyButton,
} from "@/components/ui/code-block";

export default function MariaDBSetupScript() {
  return (
    <div className="flex flex-col space-y-4">
      <CodeBlock>
        <CodeBlockHeader>
          <CodeBlockFilename value="mariadb-setup-script.sql">
            MariaDB Server Setup script
          </CodeBlockFilename>
          <CodeBlockCopyButton />
        </CodeBlockHeader>
        <CodeBlockBody value="sql">
          <CodeBlockContent language="sql">""</CodeBlockContent>
        </CodeBlockBody>
      </CodeBlock>
    </div>
  );
}
