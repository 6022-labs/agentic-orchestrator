import { Copy } from "lucide-react";
import { DropdownMenuItem } from "@/components/ui/dropdown-menu";

type CopyPairSlackUserToHumanAgentActionProps = {
  collectionTokenId: number;
  collectionAddress: string;
};

export default function CopyPairSlackUserToHumanAgentAction({
  collectionTokenId,
  collectionAddress,
}: Readonly<CopyPairSlackUserToHumanAgentActionProps>) {
  const copyToClipboard = () => {
    const text = `/pair-slack-user-to-human-agent [USER ID] ${collectionAddress} ${collectionTokenId}`;

    navigator.clipboard.writeText(text);
  };

  return (
    <DropdownMenuItem onClick={copyToClipboard}>
      <span>Copy pair slack user to human agent command</span>
      <Copy className="ml-auto h-4 w-4" />
    </DropdownMenuItem>
  );
}
