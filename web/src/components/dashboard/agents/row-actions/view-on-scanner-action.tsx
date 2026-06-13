import { ExternalLink } from "lucide-react";
import { DropdownMenuItem } from "@/components/ui/dropdown-menu";

type ViewOnScannerActionProps = {
  address: string;
};

export default function ViewOnScannerAction(
  props: Readonly<ViewOnScannerActionProps>
) {
  return (
    <DropdownMenuItem
      onClick={() => {
        window.open(
          import.meta.env.VITE_SCANNER_BASE_URL + "/address/" + props.address
        );
      }}
    >
      <span>Polygon Scan</span>
      <ExternalLink className="ml-auto h-4 w-4" />
    </DropdownMenuItem>
  );
}
