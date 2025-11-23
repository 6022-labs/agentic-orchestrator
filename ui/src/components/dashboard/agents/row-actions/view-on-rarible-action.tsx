import { ExternalLink } from "lucide-react";
import { DropdownMenuItem } from "@/components/ui/dropdown-menu";

type ViewOnRaribleActionProps = {
  collectionTokenId: number;
  collectionAddress: string;
};

export default function ViewOnRaribleAction(
  props: Readonly<ViewOnRaribleActionProps>
) {
  if (
    !import.meta.env.VITE_RARIBLE_VIEW_URL_TEMPLATE ||
    !props.collectionTokenId ||
    !props.collectionAddress
  ) {
    return null;
  }

  return (
    <DropdownMenuItem
      onClick={() => {
        window.open(
          import.meta.env.VITE_RARIBLE_VIEW_URL_TEMPLATE.replace(
            "{collection}",
            props.collectionAddress
          ).replace("{tokenId}", props.collectionTokenId),
          "_blank"
        );
      }}
    >
      <span>View on Rarible</span>
      <ExternalLink className="ml-auto h-4 w-4" />
    </DropdownMenuItem>
  );
}
