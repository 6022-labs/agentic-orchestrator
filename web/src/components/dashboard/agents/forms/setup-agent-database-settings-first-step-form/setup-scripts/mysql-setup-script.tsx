import {
  CodeBlock,
  CodeBlockBody,
  CodeBlockHeader,
  CodeBlockContent,
  CodeBlockFilename,
  CodeBlockCopyButton,
} from "@/components/ui/code-block";

export default function MySQLSetupScript() {
  return (
    <div className="flex flex-col space-y-4">
      <CodeBlock>
        <CodeBlockHeader>
          <CodeBlockFilename value="mysql-setup-script.sql">
            MySQL Server Setup script
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
