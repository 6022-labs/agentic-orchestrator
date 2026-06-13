import { Copy } from "lucide-react";
import { DropdownMenuItem } from "@/components/ui/dropdown-menu";

type CopyRemoveAgentCommandActionProps = {
  collectionTokenId: number;
  collectionAddress: string;
};

export default function CopyRemoveAgentCommandAction({
  collectionTokenId,
  collectionAddress,
}: Readonly<CopyRemoveAgentCommandActionProps>) {
  const copyToClipboard = () => {
    const text = `/remove-agent ${collectionAddress} ${collectionTokenId}`;

    navigator.clipboard.writeText(text);
  };

  return (
    <DropdownMenuItem onClick={copyToClipboard}>
      <span>Copy remove agent command</span>
      <Copy className="ml-auto h-4 w-4" />
    </DropdownMenuItem>
  );
}
