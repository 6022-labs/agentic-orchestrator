import { DropdownMenuItem } from "@/components/ui/dropdown-menu";
import { ExternalLink } from "lucide-react";

type ViewOnOpenseaActionProps = {
  collectionTokenId: number;
  collectionAddress: string;
};

export default function ViewOnOpenseaAction(
  props: Readonly<ViewOnOpenseaActionProps>
) {
  if (
    !import.meta.env.VITE_OPENSEA_VIEW_URL_TEMPLATE ||
    !props.collectionTokenId ||
    !props.collectionAddress
  ) {
    return null;
  }

  return (
    <DropdownMenuItem
      onClick={() => {
        window.open(
          import.meta.env.VITE_OPENSEA_VIEW_URL_TEMPLATE.replace(
            "{collection}",
            props.collectionAddress
          ).replace("{tokenId}", props.collectionTokenId),
          "_blank"
        );
      }}
    >
      <span>View on OpenSea</span>
      <ExternalLink className="ml-auto h-4 w-4" />
    </DropdownMenuItem>
  );
}
