import { Copy } from "lucide-react";
import { DropdownMenuItem } from "@/components/ui/dropdown-menu";

type CopyAddAgentCommandActionProps = {
  collectionTokenId: number;
  collectionAddress: string;
};

export default function CopyAddAgentCommandAction({
  collectionTokenId,
  collectionAddress,
}: Readonly<CopyAddAgentCommandActionProps>) {
  const copyToClipboard = () => {
    const text = `/add-agent ${collectionAddress} ${collectionTokenId}`;

    navigator.clipboard.writeText(text);
  };

  return (
    <DropdownMenuItem onClick={copyToClipboard}>
      <span>Copy add agent command</span>
      <Copy className="ml-auto h-4 w-4" />
    </DropdownMenuItem>
  );
}
